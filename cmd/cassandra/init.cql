-- Create keyspace
CREATE KEYSPACE IF NOT EXISTS todolist
  WITH replication = {
    'class' : 'SimpleStrategy',
    'replication_factor' : 1
  }
  AND durable_writes = true;

USE todolist;

-- Users table
CREATE TABLE IF NOT EXISTS users (
  username  text PRIMARY KEY,
  password  text,
  active    boolean
);

-- Tasks table
CREATE TABLE IF NOT EXISTS tasks (
  username      text,
  project       text,
  id            text,
  content       text,
  priority      int,
  updated_time  timestamp,
  due           timestamp,
  completed     boolean,
  PRIMARY KEY ((username), project, id)
) WITH CLUSTERING ORDER BY (project ASC, id ASC);

-- Projects table
CREATE TABLE IF NOT EXISTS projects (
  username  text,
  project   text PRIMARY KEY
);