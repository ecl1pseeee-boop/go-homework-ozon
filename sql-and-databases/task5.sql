SELECT o.item AS Название_товара, o.price AS Цена
FROM orders AS o
WHERE o.price =
      (
          SELECT MIN(orders.price)
          FROM orders
          WHERE orders.price > 10000
      )
ORDER BY o.price;

