# Coop case

## Position: Software Engineer

## General information:
For this recruitment process we would like you to solve a technical case assignment to better understand your experience and skills that we see as important for this position. The main objective for this case assignment is to understand your thought process and how you would approach such a task/project in a job setting.

## The case:
We have various output channels where we want to show data from messaging platforms. One of these
messaging platforms is the micro-blogging service Mastodon.

Given this context build an application and relevant infrastructure to stream public Mastodon posts in near real-time to end users using technologies of your choice.

Instructions:
The assignment is open-ended and fictional and you’re free to interpret it and make your own decisions and assumptions. It’s up to you how much time you spend on the assignment and how complete and comprehensive your solution is, but you should deliver something that is appropriate to be used as a basis for technical assessment for the job role level you are applying for. The cases should be presented in a manner fit for passing on knowledge to colleagues.
We ask you to present the project including the source code, architecture and thought process/decisions you took along the way during the technical interview with us.

Your case solution shall be sent back to us about one workday (24 hours) ahead of the interview so we can
prepare based on your solution. 

## Planning
- [x] Create a high-level architectural diagram
- [x] Read Mastododon api documentation
- [x] Create mastodon-to-kafka service
    - [x] Create client to interact with Mastodon API
    - [x] Figure out how to consume the api
    - [x] Create a kafka producer to use in pipeline
    - [x] Send blog messages to kafka topic
- [x] Spin up a local kafka
- [] Create frontend service
    - [] Consume kafka messages
    - [] Simple web server
    - [] Push data to the frontend (maby something like this: https://www.confluent.io/blog/webify-event-streams-using-kafka-connect-http-sink/)
- [x] Docker compose setup



## Architecture
![Landscape](/drawings/architecture.svg)