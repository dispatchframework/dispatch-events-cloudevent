# Cloudevents Event Driver

This driver enables native cloudevents (~2.0) support within Dispatch

## Installation

1. Register the cloudevents driver type (the expose option means this is a "push" driver):

    ```
    $ dispatch create eventdrivertype cloudevents dispatchframework/dispatch-events-cloudevents:kubecon-demo --expose
    Created event driver type: cloudevents
    ```

2. Create an event driver from the new type:

    ```
    $ dispatch create eventdriver cloudevents-demo --name cloudevents
    Created event driver: cloudevents-demo
    ```

3. Get the URL for the event driver.  This is the URL that the eventgrid subscriber (in Azure) will push to:

    ```
    $ dispatch get eventdriver
            NAME       |     TYPE    | STATUS | SECRETS | CONFIG |                  URL                    | REASON
    -------------------------------------------------------------------------------------------------------
      cloudevents-demo | cloudevents | READY  |         |        | https://example.com/driver/dispatch/... |
    -------------------------------------------------------------------------------------------------------
    ```

4. Test the URL (use the example cloudevent in this repo)

    ```
    $ curl -i https://example.com/driver/dispatch/... -H 'Content-Type: application/cloudevents+json' -d @event.json
    HTTP/2 200
    server: nginx/1.13.12
    date: Tue, 26 Jun 2018 17:53:54 GMT
    content-length: 0
    strict-transport-security: max-age=15724800; includeSubDomains
    ```

5. Create and Subscribe a function to the event:

    - Create a simple echo function `echo.js`:
        ```javascript
        module.exports = function (context, params) {
            params["context"] = context;
            return params
        };
        ```
    - Register the function:
        ```
        $ dispatch create function echo --image nodejs echo.js
        Created function: echo
        ```
    - Subscribe to the `word.found.noun` event:
        ```
        $ dispatch create subscription echo --event-type "word.found.noun"
        created subscription: measured-caribou-626480
        ```

6. Test the workflow:

    - Post the event to the URL (again):
        ```
        $ curl -i https://example.com/driver/dispatch/... -H 'Content-Type: application/cloudevents+json' -d @event.json
        HTTP/2 200
        server: nginx/1.13.12
        date: Tue, 26 Jun 2018 17:53:54 GMT
        content-length: 0
        strict-transport-security: max-age=15724800; includeSubDomains
        ```
    - Check that the echo function was triggered:
        ```
        $ dispatch get runs echo
           ID  | FUNCTION | STATUS |           STARTED            |           FINISHED
        ----------------------------------------------------------------------------------------
          <ID> | echo     | READY  | Tue Jun 26 11:07:40 PDT 2018 | Tue Jun 26 11:07:42 PDT 2018
        ```
    - Verify the payload in the echo result (you should see the event contents and info in the result):
        ```
        $ dispatch get runs echo <ID> --json
        ...
        ```