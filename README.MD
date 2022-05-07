<h1> How to run </h1>
<ol>
<li> Provide credentials to .env</li>
<li> <code>$ docker-compose build</code> </li>
<li> <code>$ docker-compose up</code> </li>
</ol>

---
<h2> Notes </h2>

Migrations(golang-migrate) are run automatically every time api container is running. The tests follow the same rule.</br>
Database is optional(but tests still expect it, edit dockerfile to remove test run).
If some of the db setup steps fell - error will log and app will continue running. API will work as intended only cache store will be unavailable.  
Gorm logger setuped in silent mode but all errors are logging.