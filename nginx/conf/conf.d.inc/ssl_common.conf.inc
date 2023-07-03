ssl_session_timeout 1d;
ssl_session_cache shared:TLS:10m;  # about 40000 sessions

ssl_session_tickets on;

ssl_dhparam /usr/local/nginx/conf.docker/ssl.dh/dhparams4096.pem;

ssl_protocols TLSv1.3 TLSv1.2;

ssl_prefer_server_ciphers on;

ssl_ecdh_curve auto;
ssl_ciphers "ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES256-GCM-SHA384:AES256-GCM-SHA384:!CBC";

ssl_session_ticket_key /usr/local/nginx/conf.docker/ssl.dh/ticket.key;

ssl_stapling on;
ssl_stapling_verify on;

add_header X-Content-Type-Options nosniff;
more_set_headers -t 'text/html' 'x-frame-options: SAMEORIGIN' always;
more_set_headers -t 'text/html' 'x-xss-protection: 1; mode=block';

more_set_headers "Strict-Transport-Security: max-age=63072000; includeSubDomains; preload";

ssl_early_data  on;
