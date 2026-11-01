SELECT u.name, COALESCE(SUM(o.price), 0) AS price
FROM users AS u
LEFT JOIN orders AS o ON u.id = o.user_id
GROUP BY u.name
ORDER BY price DESC
