package microdb

const (
	sqlFindCargoByID = `
	SELECT 
		cargos.tracking_id, 
		cargos.origin, 
		cargos.route_specification_id, 
		cargos.itinerary_id, 
		cargos.delivery_id, 
		RouteSpecification.id AS RouteSpecification__id, 
		RouteSpecification.created_at AS RouteSpecification__created_at, 
		RouteSpecification.updated_at AS RouteSpecification__updated_at, 
		RouteSpecification.deleted_at AS RouteSpecification__deleted_at, 
		RouteSpecification.origin AS RouteSpecification__origin, 
		RouteSpecification.destination AS RouteSpecification__destination, 
		RouteSpecification.arrival_deadline AS RouteSpecification__arrival_deadline, 
		Itinerary.id AS Itinerary__id, 
		Itinerary.created_at AS Itinerary__created_at, 
		Itinerary.updated_at AS Itinerary__updated_at, 
		Itinerary.deleted_at AS Itinerary__deleted_at, 
		Delivery.id AS Delivery__id, 
		Delivery.created_at AS Delivery__created_at, 
		Delivery.updated_at AS Delivery__updated_at, 
		Delivery.deleted_at AS Delivery__deleted_at, 
		Delivery.itinerary_id AS Delivery__itinerary_id, 
		Delivery.route_specification_id AS Delivery__route_specification_id, 
		Delivery.routing_status AS Delivery__routing_status, 
		Delivery.transport_status AS Delivery__transport_status, 
		Delivery.next_expected_activity_id AS Delivery__next_expected_activity_id, 
		Delivery.last_event_id AS Delivery__last_event_id, 
		Delivery.last_known_location AS Delivery__last_known_location, 
		Delivery.current_voyage AS Delivery__current_voyage, 
		Delivery.eta AS Delivery__eta, 
		Delivery.is_misdirected AS Delivery__is_misdirected, 
		Delivery.is_unloaded_at_destination AS Delivery__is_unloaded_at_destination 
	FROM 
		cargos 
		LEFT JOIN route_specifications RouteSpecification ON cargos.route_specification_id = RouteSpecification.id 
		LEFT JOIN itineraries Itinerary ON cargos.itinerary_id = Itinerary.id 
		LEFT JOIN deliveries Delivery ON cargos.delivery_id = Delivery.id 
	WHERE 
		tracking_id = ?`

	sqlFindAllCargos = `
	SELECT 
		cargos.tracking_id, 
		cargos.origin, 
		cargos.route_specification_id, 
		cargos.itinerary_id, 
		cargos.delivery_id, 
		RouteSpecification.id AS RouteSpecification__id, 
		RouteSpecification.created_at AS RouteSpecification__created_at, 
		RouteSpecification.updated_at AS RouteSpecification__updated_at, 
		RouteSpecification.deleted_at AS RouteSpecification__deleted_at, 
		RouteSpecification.origin AS RouteSpecification__origin, 
		RouteSpecification.destination AS RouteSpecification__destination, 
		RouteSpecification.arrival_deadline AS RouteSpecification__arrival_deadline, 
		Itinerary.id AS Itinerary__id, 
		Itinerary.created_at AS Itinerary__created_at, 
		Itinerary.updated_at AS Itinerary__updated_at, 
		Itinerary.deleted_at AS Itinerary__deleted_at, 
		Delivery.id AS Delivery__id, 
		Delivery.created_at AS Delivery__created_at, 
		Delivery.updated_at AS Delivery__updated_at, 
		Delivery.deleted_at AS Delivery__deleted_at, 
		Delivery.itinerary_id AS Delivery__itinerary_id, 
		Delivery.route_specification_id AS Delivery__route_specification_id, 
		Delivery.routing_status AS Delivery__routing_status, 
		Delivery.transport_status AS Delivery__transport_status, 
		Delivery.next_expected_activity_id AS Delivery__next_expected_activity_id, 
		Delivery.last_event_id AS Delivery__last_event_id, 
		Delivery.last_known_location AS Delivery__last_known_location, 
		Delivery.current_voyage AS Delivery__current_voyage, 
		Delivery.eta AS Delivery__eta, 
		Delivery.is_misdirected AS Delivery__is_misdirected, 
		Delivery.is_unloaded_at_destination AS Delivery__is_unloaded_at_destination 
	FROM 
		cargos 
		LEFT JOIN route_specifications RouteSpecification ON cargos.route_specification_id = RouteSpecification.id 
		LEFT JOIN itineraries Itinerary ON cargos.itinerary_id = Itinerary.id 
		LEFT JOIN deliveries Delivery ON cargos.delivery_id = Delivery.id
	LIMIT 50000`

	sqlInsertCargo = `
	INSERT INTO cargos (
		tracking_id, origin, route_specification_id, 
		itinerary_id, delivery_id
	) 
	VALUES 
		(?, ?, ?, ?, ?) ON DUPLICATE KEY 
	UPDATE 
		tracking_id = tracking_id`

	sqlFindLocationByID = `SELECT * FROM locations WHERE un_locode = ?`

	sqlFindAllLocations = `SELECT * FROM locations`

	sqlFindVoyageByID = `
	SELECT 
		voyages.number, 
		voyages.schedule_id, 
		Schedule.id AS Schedule__id, 
		Schedule.created_at AS Schedule__created_at, 
		Schedule.updated_at AS Schedule__updated_at, 
		Schedule.deleted_at AS Schedule__deleted_at 
	FROM 
		voyages 
		LEFT JOIN schedules Schedule ON voyages.schedule_id = Schedule.id
	WHERE 
		number = ?`

	sqlInsertHandlingEvent = `
	INSERT INTO handling_events (
		tracking_id, activity_id, handling_history_refer
	) 
	VALUES 
		(?, ?, ?) ON DUPLICATE KEY 
	UPDATE 
		tracking_id = tracking_id`
)
