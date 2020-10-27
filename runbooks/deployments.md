# Deployments Runbook

### My changes are not being reflected in production after a deployment

1. Check that the changes have been committed *and pushed*. The `deploy` tool uses a mirror of the repository, so all changes must exist on the remote.
2. Did the commit hash actually change? E.g. if you only made changes to the deploy tool itself, or a Dockerfile, and therefore did not need to push anything, the commit hash won't have changed. Helm won't replace the pods in this case because the resulting Docker image will have the same name and tag.
