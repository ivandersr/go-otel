# Weather Realtime Check

This app checks the realtime weather for the queried city, using its CEP.
<b>Usage:</b> Query the route with the desired CEP corresponding to the city desired by:
`https://go-weather-144114506483.us-central1.run.app/weather?cep=87225-000`

The query param needs to be 8 digits long, and may or not have a CEP mask.
Anything other than the requested format will return an error.

I used the `gcr.io/distroless/static-debian12` image as base due to problems with certificate verification when `scratch` was used. This is an image with a size about the same as scratch for static builds as this one, bringing some extra features such as certificates.
