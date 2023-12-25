# Bookings and Reservations

This is the repository for my bookings and reservations project.

- Built in Go version 1.19
- Uses [scs v2.5.1](github.com/alexedwards/scs/v2)
- Uses [chi v5.0.10](github.com/go-chi/chi/v5)
- Uses [nosurf v1.1.1](github.com/justinas/nosurf)

# Local development

1. Build image\
`docker build -t booking .`
2. Run container locally\
`docker run -p 8080:8080 -it --rm --name bookings booking`

# Deployment

Deployment using Render allows env vars to be used as build args.\
In Render settings - set env=production

1. Build image\
`docker build -t booking .`
2. Override run command with prod db\
`./bookings -dbname=bookings_ge86 -dbuser=lld -dbpass=mZS0Jcw6yAfriKCqRo7x17INg1Nvq4Pc -dbhost=dpg-cm45pqun7f5s73btffb0-a -cache=false`