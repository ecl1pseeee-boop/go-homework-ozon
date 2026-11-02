SELECT name, age
FROM users
WHERE users.city = 'Москва'
  AND age = (
        SELECT MAX(age) FROM users WHERE city = 'Москва'
    )
