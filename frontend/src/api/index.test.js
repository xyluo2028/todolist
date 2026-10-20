import axios from 'axios';
import { api } from './index';

jest.mock('axios', () => ({
  create: jest.fn(() => ({
    get: jest.fn(),
    post: jest.fn(),
    delete: jest.fn(),
  })),
}));

const client = axios.create.mock.results[0].value;

beforeEach(() => {
  jest.clearAllMocks();
});

test.each([
  ['ASCII', 'alice', 'password'],
  ['Latin-1', 'André', 'contraseña'],
  ['Chinese and Japanese', '测试用户', 'パスワード'],
  ['emoji', 'user🧑', 'secret🔐:value'],
])('%s credentials use UTF-8 Basic authentication', (_, username, password) => {
  api.getProjects({ username, password });

  expect(client.get).toHaveBeenCalledWith('/printProjects', {
    headers: {
      Authorization: `Basic ${Buffer.from(`${username}:${password}`, 'utf8').toString('base64')}`,
    },
  });
});
