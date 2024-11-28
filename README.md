# Building a Go API for Current Toronto Time with MySQL Database Logging

#### Overview 
This API will provide the current time in Toronto in JSON format. Additionally, each request to the API will log the current time to a MySQL database.

/current-time: Returns the current time in Toronto and logs it to the database.

/all-times: Retrieves and returns a list of all logged times in the database.


#### Requirements
- Go Programming Language 
-	MYSQL Database
- Go Modules
- Docker (optional for running MySQL in a container)

#### Instructions for setting up MYSQL:
- Click on the given  [Link](https://dev.mysql.com/downloads/installer/)

![Screenshot 2024-11-28 103128](https://github.com/user-attachments/assets/a8dc9ff2-5eef-44cc-9e8f-c8501e903942)

#### Script or instructions for setting up the MySQL database and table.
- Create a new database 

![image](https://github.com/user-attachments/assets/e2c0e25b-d4da-437b-b3df-ab8992f4af15)

- Create a table named time_log:

![image](https://github.com/user-attachments/assets/f980b420-17b9-4f69-b84b-81d90362462d)

#### Source code of the go application
- Installing the necessary dependencies:

![image](https://github.com/user-attachments/assets/cdff93a5-90b0-40d0-a88b-33b7ef4b3d5c)

- Start the server:

![image](https://github.com/user-attachments/assets/ba22702d-24fd-4e4c-a628-a758f2539f4c)

- Access the API

![image](https://github.com/user-attachments/assets/900e6f3b-0372-40e7-90bf-89d759c33150)

- Using Curl

![image](https://github.com/user-attachments/assets/826962f5-f724-4821-bdf6-29beb70f5e90)

- Check the database

![image](https://github.com/user-attachments/assets/586c7388-a1e5-4b68-a6f8-038b6c8db27f)

### Bonus Challenges
#### Implement logging in your Go application to log events and errors.

- Logs captured in app.log

![image](https://github.com/user-attachments/assets/f07c80c3-57ef-4638-bf60-1a75defd46cc)

![image](https://github.com/user-attachments/assets/bb182787-fd16-46be-8acf-339c1def79d7)

![image](https://github.com/user-attachments/assets/57669218-ca52-4243-92aa-bbfdca62b3d7)

#### Create an additional endpoint to retrieve all logged times from the database.

- End point "/all-timesadds" the new endpoint to handle requests to /all-times.

- Testing of the endpoint
  
Using browser

![image](https://github.com/user-attachments/assets/63ffa34f-d8e6-43e4-98cf-5d68073c7880)

Using curl

![image](https://github.com/user-attachments/assets/058566a9-623e-4072-ac18-b28d3b22c0d2)

#### Dockerize your Go application and the MySQL database for easy deployment.

- Create a Dockerfile in the root of your Go project. The Dockerfile defines the environment for your Go application and how it should be built and run in a Docker container.

![image](https://github.com/user-attachments/assets/0a9fee98-6035-4d05-8703-3533494deb27)

- Create a docker-compose.yml File: Docker Compose allows you to define and manage multi-container Docker applications. In this case, it will manage both the Go application and MySQL database.

![image](https://github.com/user-attachments/assets/bba250d7-7fff-46af-b70b-9dcdbb61b985)

#### Dockerising the application.

![image](https://github.com/user-attachments/assets/11dd7e27-f41d-4f97-bce2-1a072bdc50c5)

![image](https://github.com/user-attachments/assets/2f2969f4-d2f4-4914-a6bb-f47a3e1f4091)



















