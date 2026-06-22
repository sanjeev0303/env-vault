# Env Vault Role-Based Access Control (RBAC) Permission Matrix

The application uses a granular permission system. Roles are mapped to sets of permissions.

| Action / Permission | Super Admin | Owner | Admin | Developer | Viewer |
|--------------------|-------------|-------|-------|-----------|--------|
| Create Project     | YES         | YES*  | NO    | NO        | NO     |
| Manage Project     | YES         | YES   | NO    | NO        | NO     |
| Delete Project     | YES         | YES   | NO    | NO        | NO     |
| Add/Remove Members | YES         | YES   | YES   | NO        | NO     |
| Create Environment | YES         | YES   | YES   | NO        | NO     |
| Delete Environment | YES         | YES   | NO    | NO        | NO     |
| Create Secret      | YES         | YES   | YES   | YES       | NO     |
| View Secret (Key)  | YES         | YES   | YES   | YES       | YES    |
| Reveal Secret (Val)| YES         | YES   | YES   | YES       | NO     |
| Update Secret      | YES         | YES   | YES   | YES       | NO     |
| Delete Secret      | YES         | YES   | YES   | NO        | NO     |
| Export All Secrets | YES         | YES   | YES   | NO        | NO     |
| View Audit Logs    | YES         | YES   | YES   | NO        | NO     |

* *Global Project Creation is available to any user. The creator automatically becomes the `Owner` of that project.*
* *Super Admins bypass all project-level RBAC checks and have implicit access to everything.*
