permit(
    principal,
    action == Action::"get_inventory",
    resource == Resource::"resource-test-23"
)
when {
    principal.scopes.contains(
        "inventory:read"
    )
};

permit(
    principal,
    action == Action::"get_inventory",
    resource == Resource::"resource-test-23"
)
when {
    principal.roles.contains(
        "admin"
    )
};