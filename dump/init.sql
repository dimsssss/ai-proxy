-- `ai-proxy`.services definition

CREATE TABLE `services` (
  `id` char(36) NOT NULL DEFAULT (uuid()),
  `name` varchar(100) DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- `ai-proxy`.rate_limits definition

CREATE TABLE `rate_limits` (
  `id` char(36) NOT NULL DEFAULT (uuid()),
  `service_id` char(36) NOT NULL,
  `bpm` int NOT NULL,
  `rpm` int NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `ai-proxy`.services
(id, name, created_at)
VALUES('ac85de66-6b64-11ef-82d6-0242ac150002', 'service1', '2024-09-05 08:56:00');
INSERT INTO `ai-proxy`.services
(id, name, created_at)
VALUES('b86b7048-6b64-11ef-82d6-0242ac150002', 'service2', '2024-09-05 08:56:20');


INSERT INTO `ai-proxy`.rate_limits
(id, service_id, bpm, rpm, created_at, updated_at)
VALUES('07afdbc9-6b65-11ef-82d6-0242ac150002', 'ac85de66-6b64-11ef-82d6-0242ac150002', 10, 10, '2024-09-05 08:58:33', NULL);
INSERT INTO `ai-proxy`.rate_limits
(id, service_id, bpm, rpm, created_at, updated_at)
VALUES('f01faf15-6b64-11ef-82d6-0242ac150002', 'b86b7048-6b64-11ef-82d6-0242ac150002', 20, 20, '2024-09-05 08:57:54', NULL);