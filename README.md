# Bookings and Reservations

This is the repository for my bookings and reservations project.

- Built in Go version 1.19
- Uses [scs v2.5.1](github.com/alexedwards/scs/v2)
- Uses [chi v5.0.10](github.com/go-chi/chi/v5)
- Uses [nosurf v1.1.1](github.com/justinas/nosurf)

# Local development

1. Build image\
`docker build -t booking .`
2. Run migrations for local db (requires creation of local postgres db with config as in database.yml) \
`docker run --rm --name bookings booking soda migrate -e local`
3. Run container locally\
`docker run -p 8080:8080 -it --rm --name bookings booking`

# Deployment

Deployment on the free tier of render does not allow custom build statements, so the following settings should be added.

1. Build image\
`docker build -t booking .`
2. Custom command to run production migration\
`docker run --rm --name bookings booking soda migrate -e production`
3. Override run command with prod db\
`./bookings -dbname=bookings_ge86 -dbuser=lld -dbpass=mZS0Jcw6yAfriKCqRo7x17INg1Nvq4Pc -dbhost=dpg-cm45pqun7f5s73btffb0-a -cache=false`