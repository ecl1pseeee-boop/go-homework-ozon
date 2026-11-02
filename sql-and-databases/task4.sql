SELECT DISTINCT u.name
FROM users AS u
JOIN orders AS o ON u.id = o.user_id
WHERE o.item = 'наушники'
ORDER BY u.name;

