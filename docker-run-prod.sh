if [ "$ENV" -eq "production" ]; then
  ./bookings -dbname=bookings_ge86 -dbuser=lld -dbpass=mZS0Jcw6yAfriKCqRo7x17INg1Nvq4Pc -dbhost=dpg-cm45pqun7f5s73btffb0-a -cache=false
fi

if [ "$ENV" -eq "development" ]; then
  ./bookings -dbname=bookings -dbuser=bytedance -cache=false -production=false
fi

