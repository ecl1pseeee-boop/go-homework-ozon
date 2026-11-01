SELECT o.item AS Название_товара, o.price AS Цена
FROM users AS u
JOIN orders AS o ON u.id = o.user_id
WHERE o.price =
      (
          SELECT MIN(orders.price)
          FROM orders
          WHERE orders.price > 10000
      )
ORDER BY o.price

