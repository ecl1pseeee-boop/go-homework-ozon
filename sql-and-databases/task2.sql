SELECT city, COUNT(*) as counter
FROM users
GROUP BY city
HAVING COUNT(*) > 1
ORDER BY counter DESC, city;
