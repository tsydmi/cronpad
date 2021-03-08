
# <img src="https://github.com/ts-dmitry/cronpad/blob/master/frontend/src/assets/icon.svg" alt="Cronpad" title="Cronpad" width="40" height="40"/> Cronpad

Cronpad is a simple web time tracking application that helps to manage your and your team effectivity. 

* [Installation](#installation)
* [Usage](#usage)
* [Technologies](#technologies)
* [License](#license)

## Installation 
1. Replace `<host-address>` in *docker-compose.yml* file by your ip address or host name
2. Run using docker-compose:
```
docker-compose up -d
```
3. Open `http://localhost` in your browser

## Usage
User can have one of 3 roles:
- user (can define events)
- manager (user permissions + can define project tags + has access to project statistics)
- admin (manager permissions + can define basic tags and projects + has access to user statistics)



To change user role:
1. Open `http://localhost/auth` in your browser and login to keycloak administration console (default credentials: `admin/admin`).
2. Go to *Users -> Edit -> Role Mappings*.

## Technologies
Current configuration comprised of:
- Cronpad container with backend (Golang) and frontend(Javascript/VueJs)
- MongoDB
- Keycloak with predefined configuration (`cronpad` realm, `vue-frontend` client and 2 custom roles (`project-manager` and `admin`) to have access to list of users)  

## License
Licensed under the [GPL-3.0](LICENSE) License.
