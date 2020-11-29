db.createUser({
    user: "fluentd",
    pwd: "fluentd",
    roles: [
        {
            role: "readWrite",
            db: "home_automation_logs"
        }
    ]
});
