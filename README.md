<div align="center">
<h1>goAIRA (Automated Issue Reporting Assistnat)</h1>
 A lightweight Go backend application that automates issue reporting using exposed REST APIs. It is an ideal tool for mid-sized to large organizations where sharing API tokens or other credentials poses a security risk.
</div>

## Pre-requisits :

| Application | Resource |
|-------------|----------|
| Go | [Reference](https://go.dev/doc/install) |
| docker | [Reference](https://docs.docker.com/engine/install/) |
| Git | [Reference](https://git-scm.com/downloads) |
| Postman | [Reference](https://www.postman.com/downloads/) |

> <b>Note</b> : You can use any API testing tool instead of postman as long as you are sending required headers and body content as part of API call

## Setup :

<details>

1. Clone [goAIRA](https://github.com/rebelnato/goAIRA) git repository to local .
```bash
git clone https://github.com/rebelnato/goAIRA.git
```
2. Create `.env` file in parent directory .
> i.e : If you are cloning repository inside `path/goAIRA` then location of `.env` file should be `path/goAIRA/.env`
3. Add `Vault_pass` attribute inside `.env` file . We'll update the vault token in here so that same can be added as environment variable in docker container.
```bash
Vault_pass=<vault access token>
```
4. Create `config.yml` file in parent directory , same as `.env` file . Add servicenow endpoints and consumer ids in the yaml file.
```yaml
endpoints:
  servicenow:  
    base: "https://<your-instance-id>.service-now.com/"
  vault:
    addr1: "localhost:8300"
    addr2: "vault:8200"

consumers:
  - "Test1"
  - "Test2"
```
5. Use below command to build and run docker container utilizing config of `docker-compose.yml` file .
```bash
docker-compose up -d --build
```
`-d` : Runs conatiner in background.

`--build` : Forces docker image rebuild before container starts.

6. Run below command to modify ownership and permission for `/vault` folder .
```bash
docker exec -it vault sh -c "chown -R 100:100 /vault && chmod -R 750 /vault"
```
7. Run below command to initiate vault , which will return 5 unseal key and 1 token key . Store the keys safe as we will need the in next steps to unseal vault and login.
```bash
docker exec -it vault sh -c "chown -R 100:100 /vault && chmod -R 750 /vault"

Sample output:

Unseal Key 1: <unseal_key_1>
Unseal Key 2: <unseal_key_2>
Unseal Key 3: <unseal_key_3>
Unseal Key 4: <unseal_key_4>
Unseal Key 5: <unseal_key_5>

Initial Root Token: <Root key>
```
8. Add the `Root key` in `.env` and rebuild the containers using below commands.
```bash
docker-compose down
docker-compose up -d --build
```
><b>Note :</b> Doing this will add `Root key` as environment variable of your docker container . Same will be utilized for authentication later.
9. Run below commands to unseal vault .
```bash
docker exec -it vault vault operator unseal <unseal_key_1>
docker exec -it vault vault operator unseal <unseal_key_2>
docker exec -it vault vault operator unseal <unseal_key_3>
```
10. Goto `http://localhost:8200` and use root key capture while initiating vault to login.
11. Create a new secret engine name `secret` .
![alt text](resources/vault_secret_creation.gif)
12. Create 2 new secrets `SNOW` and `SNOW_refresh` inside `secret` engine .
![alt text](<resources/Recording 2025-03-10 231451.gif>)
13. Store below mentioned creds in associated secrets.
`SNOW` :
```json
{
  "client_id": "<ServiceNow client id>",
  "client_secret": "<Client password>",
  "password": "<ServiceNow users password>",
  "username": "<ServiceNow user name>"
}
```

`SNOW_refresh` :
```json
{
  "refresh_epoch_time": 0,
  "refresh_token": "<Keep it as "" , it will automatically be updated with appropriate vaule>"
}
```

![alt text](resources/giphy.gif)

</details>

## Exposd APIs :

> Below table contains list of exposed APIs by goAIRA program and required headers for APIs.

| Endpoint | Method | Required Headers | Use case | Example |
|----------|--------|------------------|----------|---------|
| /health | GET | No headers required | Return program health along with vault connection status | <pre><code>{<br>&nbsp;"server": "pong",<br>&nbsp;"vault": true<br>}</code></pre> |
| /createincident | POST | <pre><code>cosumerid , shortDesc , desc , caller , channel , impact , urgency</code></pre> | Can be used for creating new SNOW incident | <pre><code>{<br>&nbsp;"data": {<br>&nbsp;&nbsp;"incidentURL": "https://<your-instance-id>.service-now.com/now/nav/ui/classic/params/target/incident.do%3Fsys_id%<generated sys_id>",<br>&nbsp;&nbsp;"number": "INC0010051"<br>&nbsp;},<br>&nbsp;"status": "success"<br>}</code></pre> |
| /getincident | GET | <pre><code>consumerid , incidentNum</code></pre> | Get incident details | <pre><code>{<br>&nbsp;"result": [<br>&nbsp;&nbsp;{<br>&nbsp;&nbsp;&nbsp;"active": "true",<br>&nbsp;&nbsp;&nbsp;"activity_due": "",<br>&nbsp;&nbsp;&nbsp;"additional_assignee_list": "",<br>&nbsp;&nbsp;&nbsp;"approval": "not requested",<br>&nbsp;&nbsp;&nbsp;"approval_history": "",<br>&nbsp;&nbsp;&nbsp;"approval_set": "",<br>&nbsp;&nbsp;&nbsp;"assigned_to": "",<br>&nbsp;&nbsp;&nbsp;"assignment_group": "",<br>&nbsp;&nbsp;&nbsp;"business_duration": "",<br>&nbsp;&nbsp;&nbsp;"business_impact": "",<br>&nbsp;&nbsp;&nbsp;"business_service": "",<br>&nbsp;&nbsp;&nbsp;"business_stc": "",<br>&nbsp;&nbsp;&nbsp;"calendar_duration": "",<br>&nbsp;&nbsp;&nbsp;"calendar_stc": "",<br>&nbsp;&nbsp;&nbsp;"caller_id": {<br>&nbsp;&nbsp;&nbsp;&nbsp;"link": "https://<your-instance-id>.service-now.com/api/now/v1/table/sys_user/<sys_user>",<br>&nbsp;&nbsp;&nbsp;&nbsp;"value": "<data>"<br>&nbsp;&nbsp;&nbsp;},<br>&nbsp;&nbsp;&nbsp;"category": "inquiry",<br>&nbsp;&nbsp;&nbsp;"cause": "",<br>&nbsp;&nbsp;&nbsp;"caused_by": "",<br>&nbsp;&nbsp;&nbsp;"child_incidents": "0",<br>&nbsp;&nbsp;&nbsp;"close_code": "",<br>&nbsp;&nbsp;&nbsp;"close_notes": "",<br>&nbsp;&nbsp;&nbsp;"closed_at": "",<br>&nbsp;&nbsp;&nbsp;"closed_by": "",<br>&nbsp;&nbsp;&nbsp;"cmdb_ci": "",<br>&nbsp;&nbsp;&nbsp;"comments": "",<br>&nbsp;&nbsp;&nbsp;"comments_and_work_notes": "",<br>&nbsp;&nbsp;&nbsp;"company": "",<br>&nbsp;&nbsp;&nbsp;"contact_type": "self-service",<br>&nbsp;&nbsp;&nbsp;"contract": "",<br>&nbsp;&nbsp;&nbsp;"correlation_display": "",<br>&nbsp;&nbsp;&nbsp;"correlation_id": "",<br>&nbsp;&nbsp;&nbsp;"delivery_plan": "",<br>&nbsp;&nbsp;&nbsp;"delivery_task": "",<br>&nbsp;&nbsp;&nbsp;"description": "",<br>&nbsp;&nbsp;&nbsp;"due_date": "",<br>&nbsp;&nbsp;&nbsp;"escalation": "0",<br>&nbsp;&nbsp;&nbsp;"expected_start": "",<br>&nbsp;&nbsp;&nbsp;"follow_up": "",<br>&nbsp;&nbsp;&nbsp;"group_list": "",<br>&nbsp;&nbsp;&nbsp;"hold_reason": "",<br>&nbsp;&nbsp;&nbsp;"impact": "1",<br>&nbsp;&nbsp;&nbsp;"incident_state": "1",<br>&nbsp;&nbsp;&nbsp;"knowledge": "false",<br>&nbsp;&nbsp;&nbsp;"location": "",<br>&nbsp;&nbsp;&nbsp;"made_sla": "true",<br>&nbsp;&nbsp;&nbsp;"notify": "1",<br>&nbsp;&nbsp;&nbsp;"number": "INC0010051",<br>&nbsp;&nbsp;&nbsp;"opened_at": "2025-03-21 18:29:51",<br>&nbsp;&nbsp;&nbsp;"opened_by": {<br>&nbsp;&nbsp;&nbsp;&nbsp;"link": "https://<your-instance-id>.service-now.com/api/now/v1/table/sys_user/<user_id>",<br>&nbsp;&nbsp;&nbsp;&nbsp;"value": "<data>"<br>&nbsp;&nbsp;&nbsp;&nbsp;},<br>&nbsp;&nbsp;&nbsp;"order": "",<br>&nbsp;&nbsp;&nbsp;"origin_id": "",<br>&nbsp;&nbsp;&nbsp;"origin_table": "",<br>&nbsp;&nbsp;&nbsp;"parent": "",<br>&nbsp;&nbsp;&nbsp;"parent_incident": "",<br>&nbsp;&nbsp;&nbsp;"priority": "2",<br>&nbsp;&nbsp;&nbsp;"problem_id": "",<br>&nbsp;&nbsp;&nbsp;"reassignment_count": "0",<br>&nbsp;&nbsp;&nbsp;"reopen_count": "0",<br>&nbsp;&nbsp;&nbsp;"reopened_by": "",<br>&nbsp;&nbsp;&nbsp;"reopened_time": "",<br>&nbsp;&nbsp;&nbsp;"resolved_at": "",<br>&nbsp;&nbsp;&nbsp;"resolved_by": "",<br>&nbsp;&nbsp;&nbsp;"rfc": "",<br>&nbsp;&nbsp;&nbsp;"route_reason": "",<br>&nbsp;&nbsp;&nbsp;"service_offering": "",<br>&nbsp;&nbsp;&nbsp;"severity": "3",<br>&nbsp;&nbsp;&nbsp;"short_description": "Temp",<br>&nbsp;&nbsp;&nbsp;"sla_due": "",<br>&nbsp;&nbsp;&nbsp;"state": "1",<br>&nbsp;&nbsp;&nbsp;"subcategory": "",<br>&nbsp;&nbsp;&nbsp;"sys_class_name": "incident",<br>&nbsp;&nbsp;&nbsp;"sys_created_by": "AIRA",<br>&nbsp;&nbsp;&nbsp;"sys_created_on": "2025-03-21 18:29:51",<br>&nbsp;&nbsp;&nbsp;"sys_domain": {<br>&nbsp;&nbsp;&nbsp;&nbsp;"link": "https://<your-instance-id>.service-now.com/api/now/v1/table/sys_user_group/global",<br>&nbsp;&nbsp;&nbsp;&nbsp;"value": "global"<br>&nbsp;&nbsp;&nbsp;},<br>&nbsp;&nbsp;&nbsp;"sys_domain_path": "/",<br>&nbsp;&nbsp;&nbsp;"sys_id": "<sys_id>",<br>&nbsp;&nbsp;&nbsp;"sys_mod_count": "0",<br>&nbsp;&nbsp;&nbsp;"sys_tags": "",<br>&nbsp;&nbsp;&nbsp;"sys_updated_by": "AIRA",<br>&nbsp;&nbsp;&nbsp;"sys_updated_on": "2025-03-21 18:29:51",<br>&nbsp;&nbsp;&nbsp;"task_effective_number": "INC0010051",<br>&nbsp;&nbsp;&nbsp;"time_worked": "",<br>&nbsp;&nbsp;&nbsp;"universal_request": "",<br>&nbsp;&nbsp;&nbsp;"upon_approval": "proceed",<br>&nbsp;&nbsp;&nbsp;"upon_reject": "cancel",<br>&nbsp;&nbsp;&nbsp;"urgency": "2",<br>&nbsp;&nbsp;&nbsp;"user_input": "",<br>&nbsp;&nbsp;&nbsp;"watch_list": "",<br>&nbsp;&nbsp;&nbsp;"work_end": "",<br>&nbsp;&nbsp;&nbsp;"work_notes": "",<br>&nbsp;&nbsp;&nbsp;"work_notes_list": "",<br>&nbsp;&nbsp;&nbsp;"work_start": ""<br>&nbsp;&nbsp;&nbsp;}<br>&nbsp;&nbsp;]<br>}</code></pre> |
| /updateincident | PATCH | <pre><code>consumerid , incidentNum <br> <b>Optional headers:</b>CloseNotes (For providing resolution notes), CloseCode (Mandatory field while resolving incident , ex: "User error") , ShortDescription , Comment , WorkNote , Description ,AssignmentGroup , Status ( For updating status of the incident , available options 1-New , 2-In Progress , 6-Resoolved )</code></pre> | Can be used to update existing incident details including incident short description , description , assignment or state of incident | <pre><code>{<br><br>&nbsp;&nbsp;"result": {<br>&nbsp;&nbsp;&nbsp;"active": "true",<br>&nbsp;&nbsp;&nbsp;"activity_due": "",<br>&nbsp;&nbsp;&nbsp;"additional_assignee_list": "",<br>&nbsp;&nbsp;&nbsp;"approval": "not requested",<br>&nbsp;&nbsp;&nbsp;"approval_history": "",<br>&nbsp;&nbsp;&nbsp;"approval_set": "",<br>&nbsp;&nbsp;&nbsp;"assigned_to": "",<br>&nbsp;&nbsp;&nbsp;"assignment_group": "",<br>&nbsp;&nbsp;&nbsp;"business_duration": "1970-01-01 00:27:51",<br>&nbsp;&nbsp;&nbsp;"business_impact": "",<br>&nbsp;&nbsp;&nbsp;"business_service": "",<br>&nbsp;&nbsp;&nbsp;"business_stc": "1671",<br>&nbsp;&nbsp;&nbsp;"calendar_duration": "1970-01-01 00:27:51",<br>&nbsp;&nbsp;&nbsp;"calendar_stc": "1671",<br>&nbsp;&nbsp;&nbsp;"caller_id": {<br>&nbsp;&nbsp;&nbsp;&nbsp;"link": "https://<your-instance-id>.service-now.com/api/now/v1/table/sys_user/<sys_user>",<br>&nbsp;&nbsp;&nbsp;&nbsp;"value": "<data>"<br>&nbsp;&nbsp;&nbsp;},<br>&nbsp;&nbsp;&nbsp;"category": "inquiry",<br>&nbsp;&nbsp;&nbsp;"cause": "",<br>&nbsp;&nbsp;&nbsp;"caused_by": "",<br>&nbsp;&nbsp;&nbsp;"child_incidents": "0",<br>&nbsp;&nbsp;&nbsp;"close_code": "",<br>&nbsp;&nbsp;&nbsp;"close_notes": "",<br>&nbsp;&nbsp;&nbsp;"closed_at": "",<br>&nbsp;&nbsp;&nbsp;"closed_by": "",<br>&nbsp;&nbsp;&nbsp;"cmdb_ci": "",<br>&nbsp;&nbsp;&nbsp;"comments": "",<br>&nbsp;&nbsp;&nbsp;"comments_and_work_notes": "",<br>&nbsp;&nbsp;&nbsp;"company": "",<br>&nbsp;&nbsp;&nbsp;"contact_type": "self-service",<br>&nbsp;&nbsp;&nbsp;"contract": "",<br>&nbsp;&nbsp;&nbsp;"correlation_display": "",<br>&nbsp;&nbsp;&nbsp;"correlation_id": "",<br>&nbsp;&nbsp;&nbsp;"delivery_plan": "",<br>&nbsp;&nbsp;&nbsp;"delivery_task": "",<br>&nbsp;&nbsp;&nbsp;"description": "",<br>&nbsp;&nbsp;&nbsp;"due_date": "",<br>&nbsp;&nbsp;&nbsp;"escalation": "0",<br>&nbsp;&nbsp;&nbsp;"expected_start": "",<br>&nbsp;&nbsp;&nbsp;"follow_up": "",<br>&nbsp;&nbsp;&nbsp;"group_list": "",<br>&nbsp;&nbsp;&nbsp;"hold_reason": "",<br>&nbsp;&nbsp;&nbsp;"impact": "1",<br>&nbsp;&nbsp;&nbsp;"incident_state": "1",<br>&nbsp;&nbsp;&nbsp;"knowledge": "false",<br>&nbsp;&nbsp;&nbsp;"location": "",<br>&nbsp;&nbsp;&nbsp;"made_sla": "true",<br>&nbsp;&nbsp;&nbsp;"notify": "1",<br>&nbsp;&nbsp;&nbsp;"number": "INC0010051",<br>&nbsp;&nbsp;&nbsp;"opened_at": "2025-03-21 18:29:51",<br>&nbsp;&nbsp;&nbsp;"opened_by": {<br>&nbsp;&nbsp;&nbsp;&nbsp;"link": "https://<your-instance-id>.service-now.com/api/now/v1/table/sys_user/<sys_user>",<br>&nbsp;&nbsp;&nbsp;&nbsp;"value": "<data>"<br>&nbsp;&nbsp;&nbsp;},<br>&nbsp;&nbsp;&nbsp;"order": "",<br>&nbsp;&nbsp;&nbsp;"origin_id": "",<br>&nbsp;&nbsp;&nbsp;"origin_table": "",<br>&nbsp;&nbsp;&nbsp;"parent": "",<br>&nbsp;&nbsp;&nbsp;"parent_incident": "",<br>&nbsp;&nbsp;&nbsp;"priority": "2",<br>&nbsp;&nbsp;&nbsp;"problem_id": "",<br>&nbsp;&nbsp;&nbsp;"reassignment_count": "0",<br>&nbsp;&nbsp;&nbsp;"reopen_count": "1",<br>&nbsp;&nbsp;&nbsp;"reopened_by": {<br>&nbsp;&nbsp;&nbsp;&nbsp;"link": "https://<your-instance-id>.service-now.com/api/now/v1/table/sys_user/<sys_user>",<br>&nbsp;&nbsp;&nbsp;&nbsp;"value": "<data>"<br>&nbsp;&nbsp;&nbsp;},<br>&nbsp;&nbsp;&nbsp;"reopened_time": "2025-03-21 18:58:58",<br>&nbsp;&nbsp;&nbsp;"resolved_at": "",<br>&nbsp;&nbsp;&nbsp;"resolved_by": "",<br>&nbsp;&nbsp;&nbsp;"rfc": "",<br>&nbsp;&nbsp;&nbsp;"route_reason": "",<br>&nbsp;&nbsp;&nbsp;"service_offering": "",<br>&nbsp;&nbsp;&nbsp;"severity": "3",<br>&nbsp;&nbsp;&nbsp;"short_description": "Temp",<br>&nbsp;&nbsp;&nbsp;"sla_due": "",<br>&nbsp;&nbsp;&nbsp;"state": "1",<br>&nbsp;&nbsp;&nbsp;"subcategory": "",<br>&nbsp;&nbsp;&nbsp;"sys_class_name": "incident",<br>&nbsp;&nbsp;&nbsp;"sys_created_by": "AIRA",<br>&nbsp;&nbsp;&nbsp;"sys_created_on": "2025-03-21 18:29:51",<br>&nbsp;&nbsp;&nbsp;"sys_domain": {<br>&nbsp;&nbsp;&nbsp;&nbsp;"link": "https://<your-instance-id>.service-now.com/api/now/v1/table/sys_user_group/global",<br>&nbsp;&nbsp;&nbsp;&nbsp;"value": "global"<br>&nbsp;&nbsp;&nbsp;},<br>&nbsp;&nbsp;&nbsp;"sys_domain_path": "/",<br>&nbsp;&nbsp;&nbsp;"sys_id": "<sys_id>",<br>&nbsp;&nbsp;&nbsp;"sys_mod_count": "5",<br>&nbsp;&nbsp;&nbsp;"sys_tags": "",<br>&nbsp;&nbsp;&nbsp;"sys_updated_by": "AIRA",<br>&nbsp;&nbsp;&nbsp;"sys_updated_on": "2025-03-21 18:58:58",<br>&nbsp;&nbsp;&nbsp;"task_effective_number": "INC0010051",<br>&nbsp;&nbsp;&nbsp;"time_worked": "",<br>&nbsp;&nbsp;&nbsp;"universal_request": "",<br>&nbsp;&nbsp;&nbsp;"upon_approval": "proceed",<br>&nbsp;&nbsp;&nbsp;"upon_reject": "cancel",<br>&nbsp;&nbsp;&nbsp;"urgency": "2",<br>&nbsp;&nbsp;&nbsp;"user_input": "",<br>&nbsp;&nbsp;&nbsp;"watch_list": "",<br>&nbsp;&nbsp;&nbsp;"work_end": "",<br>&nbsp;&nbsp;&nbsp;"work_notes": "",<br>&nbsp;&nbsp;&nbsp;"work_notes_list": "",<br>&nbsp;&nbsp;&nbsp;"work_start": ""<br>&nbsp;&nbsp;}<br>}</code></pre> |


## Utilization / API consumption :
