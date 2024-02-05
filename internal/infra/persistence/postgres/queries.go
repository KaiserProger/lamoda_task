package postgres

const (
	getStoredItemsQuery = `
		SELECT 
			warehouses_items.item_code,
			warehouses_items.warehouse_id,
			warehouses_items.quantity
		FROM
			warehouses_items
		LEFT JOIN
			warehouses ON warehouses_items.warehouse_id = warehouses.id 
		WHERE
			warehouses.accessible = true
			AND warehouses_items.item_code = ANY ($1::int[])
			AND warehouses_items.quantity > 0
		ORDER BY warehouses_items.item_code;`
	getWarehouseBaseQuery = `
		SELECT *
		FROM warehouses
		WHERE id = $1;`
	getWarehouseItemsQuery = `
		SELECT
			items.code,
			items.name,
			items.size,
			warehouses_items.quantity
		FROM warehouses_items
		LEFT JOIN items ON warehouses_items.item_code = items.code
		WHERE warehouses_items.warehouse_id = $1;
	`
	getReservationsQuery = `
		SELECT item_code, warehouse_id, quantity FROM reservations
		WHERE item_code = ANY ($1::int[])
		ORDER BY created_at DESC;`
	makeReservationQuery = `
		INSERT INTO reservations (item_code, warehouse_id, quantity)
		VALUES ?
		ON CONFLICT (item_code, warehouse_id) DO
		UPDATE SET quantity=reservations.quantity+EXCLUDED.quantity;`
	removeFromStockQuery = `
		INSERT INTO warehouses_items (item_code, warehouse_id, quantity)
		VALUES ?
		ON CONFLICT (item_code, warehouse_id) DO
		UPDATE SET quantity=warehouses_items.quantity-EXCLUDED.quantity;`
	dereserveQuery = `
		INSERT INTO reservations (item_code, warehouse_id, quantity)
		VALUES ?
		ON CONFLICT (item_code, warehouse_id) DO
		UPDATE SET quantity=reservations.quantity-EXCLUDED.quantity
		RETURNING warehouse_id;`
	addToStockQuery = `
		INSERT INTO warehouses_items (item_code, warehouse_id, quantity)
		VALUES ?
		ON CONFLICT (item_code, warehouse_id) DO
		UPDATE SET quantity=warehouses_items.quantity+EXCLUDED.quantity;`
)
