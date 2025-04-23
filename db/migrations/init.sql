-- Create the blocked_numbers table
CREATE TABLE IF NOT EXISTS blocked_numbers (
    id SERIAL PRIMARY KEY,
    number VARCHAR(20) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Function to add a blocked number
CREATE OR REPLACE FUNCTION add_blocked_number(
    p_number VARCHAR(20),
    p_description TEXT
)
RETURNS INTEGER AS $$
DECLARE
    new_id INTEGER;
BEGIN
    INSERT INTO blocked_numbers (number, description)
    VALUES (p_number, p_description)
    RETURNING id INTO new_id;
    
    RETURN new_id;
EXCEPTION
    WHEN unique_violation THEN
        RAISE NOTICE 'Phone number already exists';
        RETURN 0;
END;
$$ LANGUAGE plpgsql;

-- Function to check if a number is blocked
CREATE OR REPLACE FUNCTION is_number_blocked(
    p_number VARCHAR(20)
)
RETURNS BOOLEAN AS $$
DECLARE
    found BOOLEAN;
BEGIN
    SELECT EXISTS(
        SELECT 1 FROM blocked_numbers 
        WHERE number = p_number
    ) INTO found;
    
    RETURN found;
END;
$$ LANGUAGE plpgsql;

-- Function to remove a blocked number
CREATE OR REPLACE FUNCTION remove_blocked_number(
    p_id INTEGER
)
RETURNS BOOLEAN AS $$
DECLARE
    rows_affected INTEGER;
BEGIN
    DELETE FROM blocked_numbers
    WHERE id = p_id
    RETURNING 1 INTO rows_affected;
    
    RETURN FOUND;
END;
$$ LANGUAGE plpgsql;

-- Function to get all blocked numbers
CREATE OR REPLACE FUNCTION get_all_blocked_numbers()
RETURNS TABLE (
    id INTEGER,
    number VARCHAR(20),
    description TEXT,
    created_at TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
    SELECT bn.id, bn.number, bn.description, bn.created_at
    FROM blocked_numbers bn
    ORDER BY bn.created_at DESC;
END;
$$ LANGUAGE plpgsql;
