-- Drop trigger first
DROP TRIGGER IF EXISTS emp_id_trigger ON employees;

-- Drop trigger function
DROP FUNCTION IF EXISTS set_emp_id();

-- Drop UID generator function
DROP FUNCTION IF EXISTS generate_uid(TEXT, TEXT);

-- Drop the employees table
DROP TABLE IF EXISTS employees;