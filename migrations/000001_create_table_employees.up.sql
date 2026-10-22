CREATE TABLE employees (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    full_name varchar(255) NOT NULL,
    position varchar(255) NOT NULL
);