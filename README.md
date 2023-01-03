## Methods

- getAllRunningTasks
- run task by id
- new task

# 1. Create new task flow

- create TaskChedule entity in db with started time [done]
- start fetch job for every locality
- persist result in DB for every page fetched with raw data (rawPages)

# 2. Create DB model for Task

# Models

--- TaskSchedule

- id
- startTime
- endTime
- locationsInQueue

--- LocationResult

- id
- locality {
  id: number;
  title: string;
  }
- taskID
- lastPage 3
- completed
- rawPages: []{url: string; html: string}
- entries: [
  {
  raw: string?
  url: string?
  attributes: {
  title: string
  price: string
  size: number
  }
  }
  ]

# 3.
