import { act } from 'react';
import { createRoot } from 'react-dom/client';
import { Simulate } from 'react-dom/test-utils';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { api } from '../api';
import Tasks from './Tasks';

jest.mock('../api', () => ({
  api: {
    getTasks: jest.fn(),
    writeTask: jest.fn(),
    completeTask: jest.fn(),
    removeTask: jest.fn(),
  },
}));

const savedTask = {
  id: 'saved-task',
  title: 'Existing task',
  content: 'Keep this task visible',
  due: '2026-10-21T10:00:00Z',
  priority: 1,
  completed: false,
};

let container;
let root;

beforeEach(() => {
  jest.resetAllMocks();
  jest.spyOn(console, 'error').mockImplementation(() => {});
  jest.spyOn(window, 'alert').mockImplementation(() => {});
  sessionStorage.setItem('auth', JSON.stringify({ username: 'alice', password: 'password' }));
  api.getTasks.mockResolvedValue({ data: [savedTask] });
  container = document.createElement('div');
  document.body.appendChild(container);
  root = createRoot(container);
});

afterEach(async () => {
  await act(async () => root.unmount());
  container.remove();
  sessionStorage.clear();
  jest.restoreAllMocks();
});

const renderTasks = async (project = 'home') => {
  await act(async () => {
    root.render(
      <MemoryRouter initialEntries={[`/projects/${encodeURIComponent(project)}`]}>
        <Routes>
          <Route path="/projects/:projectName" element={<Tasks />} />
        </Routes>
      </MemoryRouter>
    );
  });
};

const change = (selector, value) => {
  const input = container.querySelector(selector);
  act(() => Simulate.change(input, { target: { value } }));
};

const click = async (label) => {
  const button = Array.from(container.querySelectorAll('button'))
    .find((item) => item.textContent === label);
  expect(button).toBeDefined();
  await act(async () => Simulate.click(button));
};

const fillNewTask = () => {
  change('.task-inputs input[placeholder="Task title"]', 'Draft task');
  change('.task-inputs textarea', 'Unsaved details');
  change('.task-inputs input[type="datetime-local"]', '2026-10-22T12:30');
};

const submitNewTask = async () => {
  await act(async () => Simulate.submit(container.querySelector('.task-inputs')));
};

test.each(['50% done', 'Literal %20 name', '日本語 🌸 %'])(
  'opens the exact project name %s',
  async (project) => {
    await renderTasks(project);

    expect(container.querySelector('h2').textContent).toBe(`Project · ${project}`);
    expect(api.getTasks).toHaveBeenCalledWith(project, {
      username: 'alice',
      password: 'password',
    });
  }
);

test('a failed add keeps the draft and existing tasks until a successful retry', async () => {
  await renderTasks();
  fillNewTask();
  api.writeTask.mockRejectedValueOnce(new Error('Network error'));

  await submitNewTask();

  expect(container.querySelector('.task-inputs input').value).toBe('Draft task');
  expect(container.querySelector('.task-inputs textarea').value).toBe('Unsaved details');
  expect(container.querySelector('input[type="datetime-local"]').value).toBe('2026-10-22T12:30');
  expect(container.textContent).toContain('Existing task');
  expect(api.getTasks).toHaveBeenCalledTimes(1);
  expect(window.alert).toHaveBeenCalledWith(expect.stringContaining('Unable to add'));

  api.writeTask.mockResolvedValueOnce({});
  api.getTasks.mockResolvedValueOnce({ data: [savedTask, { ...savedTask, id: 'new-task', title: 'Draft task' }] });
  await submitNewTask();

  expect(container.querySelector('.task-inputs input').value).toBe('');
  expect(container.querySelectorAll('.task-card')).toHaveLength(2);
});

test('a failed edit keeps the editor open with the unsaved changes', async () => {
  await renderTasks();
  await click('Edit');
  change('.task-card--editing input[placeholder="Task title"]', 'Edited draft');
  change('.task-card--editing textarea', 'Edited details');
  api.writeTask.mockRejectedValueOnce(new Error('Server error'));

  await act(async () => Simulate.submit(container.querySelector('.task-card__edit-fields')));

  expect(api.writeTask).toHaveBeenCalledWith('home', expect.objectContaining({ title: 'Edited draft' }), expect.any(Object));
  expect(container.querySelector('.task-card--editing input').value).toBe('Edited draft');
  expect(container.querySelector('.task-card--editing textarea').value).toBe('Edited details');
  expect(container.querySelectorAll('.task-card')).toHaveLength(1);
  expect(api.getTasks).toHaveBeenCalledTimes(1);
  expect(window.alert).toHaveBeenCalledWith(expect.stringContaining('Unable to update'));
});

test.each([
  ['Mark done', 'completeTask', 'Unable to complete'],
  ['Delete', 'removeTask', 'Unable to delete'],
])('a failed %s keeps the saved list visible and reports the error', async (label, method, message) => {
  await renderTasks();
  api[method].mockRejectedValueOnce(new Error('Unauthorized'));

  await click(label);

  expect(container.textContent).toContain('Existing task');
  expect(container.querySelector('.task-card--completed')).toBeNull();
  expect(api.getTasks).toHaveBeenCalledTimes(1);
  expect(window.alert).toHaveBeenCalledWith(expect.stringContaining(message));
});

test('a failed refresh after a successful save keeps the list and retries without saving twice', async () => {
  await renderTasks();
  fillNewTask();
  api.writeTask.mockResolvedValueOnce({});
  api.getTasks.mockRejectedValueOnce(new Error('Refresh unavailable'));

  await submitNewTask();

  expect(container.textContent).toContain('Existing task');
  expect(container.querySelector('.task-inputs input').value).toBe('');
  expect(container.querySelector('[role="alert"]').textContent).toContain('Unable to refresh tasks');
  expect(window.alert).not.toHaveBeenCalled();

  api.getTasks.mockResolvedValueOnce({ data: [savedTask, { ...savedTask, id: 'new-task', title: 'Draft task' }] });
  await click('Retry');

  expect(container.querySelector('[role="alert"]')).toBeNull();
  expect(container.querySelectorAll('.task-card')).toHaveLength(2);
  expect(api.writeTask).toHaveBeenCalledTimes(1);
});
