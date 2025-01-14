# ClockifyToJira

### What is clockify-to-jira

Clockify-to-jira is an application which allows you to migrate clockify time entries into multiple jira instances.

### How it works

Clockify-to-jira takes clockify time entries from the given time period and it creates workloads in related jira issues. To work properly each toogle timeentry has to be assigned to the client. Base on this assignment and app configuration, application recognize in which jira instance should create workloads.

### Installation

MacOS Homebrew:

```
brew install kruc/homebrew-tap/clockify-to-jira
```

Linux

```
 wget -qO- https://github.com/kruc/clockify-to-jira/releases/download/v1.0.0/clockify-to-jira_1.0.0_Linux_x86_64.tar.gz | tar -xvz -C /usr/local/bin && chmod +x /usr/local/bin/clockify-to-jira
```

### Requirements

1. Clockify time entries naming convention

   ```
   [JIRA-ISSUE-ID] [WORKLOAD DESCRIPTION]
   e.g. ISSUE-123 Description of what has been done

   https://jira.atlassian.net/browse/ISSUE-123
   ```

   [JIRA-ISSUE-ID] - issue matching

   [WORKLOAD DESCRIPTION] - jira workload description

2. Assign client to every time entry you want to migrate

   jira instance matching is based on client

### First run

1. Create config file (by default `$HOME/.clockify-to-jira/config.yaml`)

   ```yaml
   global:
     clockify_token: clockify-token
     period: 1

   default_client:
     jira_client_user: firstname.lastname
     jira_host: https://jira.atlassian.net
     jira_password: (visit https://id.atlassian.com/manage/api-tokens)
     jira_username: firstname.lastname@domain.com
     stachursky_mode: 1

   default_workspace:
     jira_migration_failed_tag: jira-migration-failed
     jira_migration_skip_tag: jira-migration-skip
     jira_migration_success_tag: logged

   workspaces:
     ws_1:
       workspace_id: ws-1
       jira_migration_failed_tag: jira-migration-failed
       jira_migration_success_tag: logged
       clients:
         client_1:
           enabled: true
           jira_client_user: username
           jira_host: https://domain.atlassian.net
           jira_password: jirapassword-client-1
           jira_username: username@domain.com
           stachursky_mode: 30
         client_2:
           enabled: false

     ws_2:
       workspace_id: ws-2
       jira_migration_failed_tag: failed
       jira_migration_skip_tag: skipped
       jira_migration_success_tag: logged
       clients:
         client_3:
           enabled: false
           jira_password: jirapassword-client-3
           jira_username: username3@domain.com
   ```

1. Adjust the configuration to your needs :sweat_smile:

1. Run help command to check available options

   ```bash
   clockify-to-jira -h
   ```

1. Run migration in dry-run mode (without -a | --apply flag)

   ```bash
   clockify-to-jira -p 3 # show time entries from last 3 days
   ```

1. If everything is correct, run with the `--apply` flag

   ```bash
   clockify-to-jira -p 3 --apply
   ```

   ```bash
    [22:44:01.928] INFO: Worklog: XYZ-1445 Some description
    ---------
    Workspace: workspace
    Client: client
    Project: project
    Date: 2025-01-13 10:28:13
    Time spent: 5h30m0s (clockify: 5h27m0s stachurskyMode: 15m)
    Comment: Some comment
    Tags: []
    ---------

    [22:44:01.928] INFO: Workspace: workspace
    -------
    SUMMARY
    -------
    Time entries range: 2025-01-10 22:44:01 - 2025-01-13 22:44:01
    Number of time entries: 2
    Total time: 8h55m0s
    Total dosko: 9h0m0s (t=15m)
   ```

1. After migration success clockify time entry will be tag with `jira_migration_success_tag` configuration key value (default: `logged`) - this tag causes skip on next migration
1. If you want to skip some time entry migration, tag it with `jira_migration_skip_tag` configuration key value (default: `jira-migration-skip`)
1. After migration fail clockify time entry will be tag with `jira_migration_failed_tag` configuration key value (default: `jira-migration-failed`) - this tag will be remove after migration success
