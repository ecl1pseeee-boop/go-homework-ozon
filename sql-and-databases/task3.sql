SELECT u.name, COALESCE(SUM(o.price), 0) AS total
FROM users AS u
LEFT JOIN orders AS o ON u.id = o.user_id
GROUP BY u.id, u.name
ORDER BY total DESC, u.name;
