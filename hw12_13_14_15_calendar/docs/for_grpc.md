/CreateEvent
/UpdateEvent
```text
{
    "event": {
        "id": "1234567896",
        "title": "Meet",
        "description": "Team meet",
        "user_id": "user126",
        "duration": "120",
        "reminder": "15",
        "date_time": {
            "seconds": "1696779648",
            "nanos": 0
        }
    }
}
```

/GetEvent
/DeleteEvent
```text
{
    "event_id": "1234567892"
}
```

/ListEventsForDay
/ListEventsForWeek
/ListEventsForMonth
```text
{
  "date": {
    "seconds": 1697133416,
    "nanos": 999999999
  }
}
```
