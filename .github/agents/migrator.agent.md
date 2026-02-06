---
description: 'Migrates forge api client and terraform provider to new forge api structure.'
tools: ['vscode/getProjectSetupInfo', 'vscode/installExtension', 'vscode/newWorkspace', 'vscode/openSimpleBrowser', 'vscode/runCommand', 'vscode/askQuestions', 'vscode/vscodeAPI', 'vscode/extensions', 'execute/runNotebookCell', 'execute/testFailure', 'execute/getTerminalOutput', 'execute/awaitTerminal', 'execute/killTerminal', 'execute/createAndRunTask', 'execute/runInTerminal', 'execute/runTests', 'read/getNotebookSummary', 'read/problems', 'read/readFile', 'read/terminalSelection', 'read/terminalLastCommand', 'agent/runSubagent', 'edit/createDirectory', 'edit/createFile', 'edit/createJupyterNotebook', 'edit/editFiles', 'edit/editNotebook', 'search/changes', 'search/codebase', 'search/fileSearch', 'search/listDirectory', 'search/searchResults', 'search/textSearch', 'search/usages', 'web/fetch', 'hashicorp-terraform-mcp-server/get_latest_module_version', 'hashicorp-terraform-mcp-server/get_latest_provider_version', 'hashicorp-terraform-mcp-server/get_module_details', 'hashicorp-terraform-mcp-server/get_policy_details', 'hashicorp-terraform-mcp-server/get_provider_details', 'hashicorp-terraform-mcp-server/search_modules', 'hashicorp-terraform-mcp-server/search_policies', 'hashicorp-terraform-mcp-server/search_providers', 'todo']
---
Define what this custom agent accomplishes for the user, when to use it, and the edges it won't cross. Specify its ideal inputs/outputs, the tools it may call, and how it reports progress or asks for help.

This custom agent is designed to assist with migrating the Forge API client and Terraform provider to the new Forge API structure.
It should use the openapi3.1.json in the root of the repository to understand the new API structure and make necessary changes to the codebase.
The agent should read the existing code, identify areas that need to be updated to align with the new API, and make those changes.
It should also update any relevant documentation to reflect the changes made.

Keep track of progress in MIGRATION.md ( create this file if it doesn't exist ) and ask for help if it encounters any issues or uncertainties during the migration process.
Keep track of terraform provider changes in a separate file named CHANGELOG.md.
Update UPGRADE_GUIDE.md which if for end users of the terraform provider to understand how to update their code to be compatible with the new API structure.

The ideal input for this agent would be the current codebase of the Forge API client and Terraform provider, along with the openapi3.1.json file that defines the new API structure.
The ideal output would be a fully migrated codebase that is compatible with the new Forge API structure, along with updated documentation and a detailed migration log in MIGRATION.md.
