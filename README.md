# BCC University

### ⚠️⚠️⚠️

```
Submissions from 2024 students will have much higher priority than submissions from 2023, SAP, or higher students.
Please take note of this before planning to attempt this freepass challenge.
```

## 💌 Invitation Letter

Embracing the ever-evolving conference management landscape, we recognize the need for a seamless and engaging experience in academic meetings especially in the organizational space. We are embarking on an innovative project to transform the way conferences are hosted and experienced, and we want you to be a part of this journey!

We aim to create a dynamic conference platform that revolutionizes session management, attendee engagement, and administrative oversight. Your contributions will help shape the future of conference management. Together, we can create a platform that enhances knowledge sharing and professional networking while maintaining the highest standards of academic discourse.

Join us in revolutionizing the conference experience. Your insights and expertise are key to making this transformation happen!

## **⭐** Minimum Viable Product (MVP)

As we have mentioned earlier, we need technology that can support BCC Conference in the future. Please consider these features below:

- New user can register account to the system ✔️
- User can login to the system ✔️
- User can edit their profile account ✔️
- User can view all conference sessions ✔️
- User can leave feedback on sessions ✔️
- User can view other user's profile ✔️
- Users can register for sessions during the conference registration period if seats are available ✔️
- Users can only register for one session within a time period ✔️
- Users can create, edit, delete their session proposals ✔️
- Users can only create one session proposal within a time period ✔️
- Users can edit, delete their session ✔️
- Event Coordinator can view all session proposals ✔️
- Event Coordinator can accept or reject user session proposals ✔️
- Event Coordinator can remove sessions ✔️
- Event Coordinator can remove user feedback ✔️
- Admin can add new event coordinators ✔️
- Admin can remove users/event coordinators ✔️

## **🌎** Service Implementation

```
GIVEN => I am a new user
WHEN  => I register to the system
THEN  => System will record and return the user's registration details

GIVEN => I am a user
WHEN  => I log in to the system
THEN  => System will authenticate and grant access based on user credentials

GIVEN => I am a user
WHEN  => I edit my profile account
THEN  => The system will update my account with the new information

GIVEN => I am a user
WHEN  => I view all available conference's sessions
THEN  => System will display all conference sessions with their details

GIVEN => I am a user
WHEN  => I leave feedback on a session
THEN  => System will record my feedback and display it under the session

GIVEN => I am a user
WHEN  => I view other user's profiles
THEN  => System will show the information of other user's profiles

GIVEN => I am a user
WHEN  => I register for conference's sessions
THEN  => System will confirm my registration for the selected session

GIVEN => I am a user
WHEN  => I create a new session proposal
THEN  => System will record and confirm the session creation

GIVEN => I am a user
WHEN => I see my session's proposal details
THEN => System will display my session's proposal details

GIVEN => I am a user
WHEN  => I update my session's proposal details
THEN  => System will apply the changes and confirm the update

GIVEN => I am a user
WHEN  => I delete my session's proposal
THEN  => System will remove the session's proposal

GIVEN => I am a user
WHEN => I see my session details
THEN => System will display my session details

GIVEN => I am a user
WHEN  => I update my session details
THEN  => System will apply the changes and confirm the update

GIVEN => I am a user
WHEN  => I delete my session
THEN  => System will remove the session

GIVEN => I am an event coordinator
WHEN  => I view session proposals
THEN  => System will display all submitted session proposals

GIVEN => I am an event coordinator
WHEN  => I accept or reject the session proposal
THEN  => The system will make the session either be available or unavailable

GIVEN => I am an event coordinator
WHEN  => I remove a session
THEN  => System will delete the session

GIVEN => I am an event coordinator
WHEN  => I remove user feedback
THEN  => System will delete the feedback from the session

GIVEN => I am an admin
WHEN  => I add new event coordinator
THEN  => System will make the account to the system

GIVEN => I am an admin
WHEN  => I remove a user or event coordinator
THEN  => System will delete the account from the system
```

## **👪** Entities and Actors

We want to see your perspective about these problems. You can define various types of entities or actors. One thing for sure, there is no
true or false statement to define the entities. As long as the results are understandable, then go for it! 🚀

## **📘** References

You might be overwhelmed by these requirements. Don't worry, here's a list of some tools that you could use (it's not required to use all of them nor any of them):

1. [Example Project](https://github.com/meong1234/fintech)
2. [Git](https://try.github.io/)
3. [Cheatsheets](https://devhints.io/)
4. [REST API](https://restfulapi.net/)
5. [Insomnia REST Client](https://insomnia.rest/)
6. [Test-Driven Development](https://www.freecodecamp.org/news/test-driven-development-what-it-is-and-what-it-is-not-41fa6bca02a2/)
7. [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
8. [GraphQL](https://graphql.org/)
9. [gRPC](https://grpc.io/)
10. [Docker Compose](https://docs.docker.com/compose/install/)

## **🔪** Accepted Weapons

> BEFORE CHOOSING YOUR LANGUAGE, PLEASE VISIT OUR [CONVENTION](CONVENTION.md) ON THIS PROJECT
>
> **Any code that did not follow the convention will be rejected!**
>
> 1. Golang (preferred)
> 2. Java (preferred)
> 3. NodeJS
> 4. PHP

You are welcome to use any libraries or frameworks, but we appreciate if you use the popular ones.

## **🎒** Tasks

```
The implementation of this project MUST be in the form of a REST, gRPC, or GraphQL API (choose AT LEAST one type).
```

1. Fork this repository
2. Follow the project convention
3. Finish all service implementations
4. Write the installation guide of your back-end service in the section below

## **🧪** API Installation

> Write how to run your service in local or development environment here. If you use Docker to serve your DBMS or your server, you will receive bonus points for your submission.

> See the detail of the project [here](ABOUT.md) here.

### Run in local environment

> **Note:** This project uses [Task](https://taskfile.dev/) as a task runner. You can also run the commands manually if you don't want to use Task. Just copy the command from the `Taskfile.yml` file.

1. Ensure you have [Go](https://go.dev/dl/) 1.23 or higher installed on your machine:

   ```bash
   go version && task --version # windows
   go version && go-task --version # unix
   ```

2. Create a copy of the `.env.example` file and rename it to `.env`:

   ```bash
   cp ./config/.env.example ./config/.env
   ```

   Update configuration values as needed.

3. Install all dependencies:

   ```bash
   task
   ```

4. Initialize `air` configuration:

	 ```bash
	 air init

	 # Update the `air.toml` file with the following configuration:
	 # [build]
	 # cmd = "go build -o ./tmp/main ./cmd/app" # for unix
	 # cmd = "go build -o ./tmp/main.exe ./cmd/app" # for windows
	 ```

5. Run migrations:

	 ```bash
	 task migrate:up
	 ```

6. Seed the database:

	 ```bash
	 task db:seed
	 ```

7. Run the project in development mode:

   ```bash
   task dev
   ```

### Run in Docker environment
1. Ensure you have [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) installed on your machine:

   ```bash
   docker --version && docker-compose --version
   ```

2. Create a copy of the `.env.example` file and rename it to `.env`:

   ```bash
    cp ./config/.env.example ./config/.env
    ```

    Update configuration values as needed. Ensure the `DB_HOST` is set to `db`.

3. Build and run the project:

   ```bash
   go-task service:run -- dev
   ```

4. Run migrations:

   ```bash
   go-task service:db:migrate
   ```

   Ensure the `DB_HOST` is set to `localhost`.

5. Seed the database:

   ```bash
    go-task db:seed
    ```

    Ensure the `DB_HOST` is set to `localhost`.

> **Note:** If you encounter host resolution issues, try to update the `DB_HOST` value to `localhost` or `db` in the `.env` file.

## **📞** Contact

Have any questions? You can contact either [Tyo](https://www.instagram.com/nandanatyo/) or [Ilham](https://www.instagram.com/iilham_akbar/).

## **🎁** Submission

Please follow the instructions on the [Contributing guide](CONTRIBUTING.md).

![cheers](https://gifsec.com/wp-content/uploads/2022/10/cheers-gif-11.gif)

> This is not the only way to join us.
>
> **But, this is the _one and only way_ to instantly pass.**
