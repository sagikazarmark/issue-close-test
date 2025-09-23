permit(
	principal in Role::"Collaborator",
	action in [Action::"close", Action::"label", Action::"remove-label"],
	resource
);
