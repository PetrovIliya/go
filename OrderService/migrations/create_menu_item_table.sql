CREATE TABLE `orderservice`.`menu_item` (
    `menu_item_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(100) NOT NULL,
    `order_id` INT UNSIGNED NOT NULL,
    `quantity` INT UNSIGNED NOT NULL,
    PRIMARY KEY (`menu_item_id`)
    );