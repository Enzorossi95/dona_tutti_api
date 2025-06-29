-- +goose Up
-- Refactor donation payment method from enum to foreign key

-- 1. Add new payment_method_id column
ALTER TABLE donations ADD COLUMN payment_method_id INTEGER;

-- 2. Create foreign key constraint
ALTER TABLE donations ADD CONSTRAINT fk_donations_payment_method 
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id);

-- 3. Migrate existing data
UPDATE donations SET payment_method_id = (
    SELECT pm.id FROM payment_methods pm 
    WHERE LOWER(pm.code) = LOWER(CASE 
        WHEN donations.payment_method = 'MercadoPago' THEN 'mercadopago'
        WHEN donations.payment_method = 'Transferencia' THEN 'transfer'  
        WHEN donations.payment_method = 'Efectivo' THEN 'cash'
        ELSE 'transfer'
    END)
    LIMIT 1
);

-- 4. Make payment_method_id NOT NULL after migration
ALTER TABLE donations ALTER COLUMN payment_method_id SET NOT NULL;

-- 5. Drop the old payment_method column (enum)
ALTER TABLE donations DROP COLUMN payment_method;

-- +goose Down
-- Reverse the changes

-- 1. Add back the old payment_method column
ALTER TABLE donations ADD COLUMN payment_method VARCHAR(20);

-- 2. Migrate data back
UPDATE donations SET payment_method = (
    SELECT CASE 
        WHEN pm.code = 'mercadopago' THEN 'MercadoPago'
        WHEN pm.code = 'transfer' THEN 'Transferencia'
        WHEN pm.code = 'cash' THEN 'Efectivo'
        ELSE 'Transferencia'
    END
    FROM payment_methods pm 
    WHERE pm.id = donations.payment_method_id
);

-- 3. Make payment_method NOT NULL and drop foreign key
ALTER TABLE donations ALTER COLUMN payment_method SET NOT NULL;
ALTER TABLE donations DROP CONSTRAINT fk_donations_payment_method;
ALTER TABLE donations DROP COLUMN payment_method_id;