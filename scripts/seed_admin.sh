#!/bin/bash
# Seed script for creating initial admin user
# This is a template - implement based on your needs

echo "Note: This is a template script."
echo "To create an admin user, you'll need to:"
echo "1. Hash the password using bcrypt"
echo "2. Insert into users table with role='admin'"
echo ""
echo "Example SQL:"
echo "INSERT INTO users (id, email, password_hash, role, full_name, is_active, created_at, updated_at)"
echo "VALUES (gen_random_uuid(), 'admin@touros.gov.np', '<bcrypt_hash>', 'admin', 'System Admin', true, NOW(), NOW());"
echo ""
echo "You can use Go code or a tool like https://bcrypt-generator.com/ to generate the hash"

