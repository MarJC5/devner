#!/bin/bash

handle_credit() {
# ASCII Art
ascii_art=$(cat <<'EOF'

                           ╔@@@@@@@@@@@@_       ╔╔╔╔╔╔╔╔╔╔╔╔_
                          £╠╠╠╠╠╠╠╠╠╠╠╠╠╠²       ╚ÜÜÜÜÜÜÜÜÜÜÜ_
                         ╬╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠,       ╚ÜÜÜÜÜÜÜÜÜÜÜ╓
                        ╬╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╕       ╚ÜÜÜÜÜÜÜÜÜÜÜ╓
                       ╬╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╦       \ÜÜÜÜÜÜÜÜÜ░╠╦
                      ╬╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╠╦       ²ÜÜÜÜÜÜÜÜ╠╠╠▒
                    '╠╠╠╠╠╠╠╠╠╠╠╝ ╙╠╠╠╠╠╠╠╠╠╠╠▒       'ÜÜÜÜÜÜ╠╠╠╠╠▒
                   ╓╠╠╠╠╠╠╠╠╠╠╠╩   └╠╠╠╠╠╠╠╠╠╠╠▒       `ÜÜÜ▒╠╠╠╠╠╠╠▒
                  /╠╠╠╠╠╠╠╠╠╠╠╩     '╠╠╠╠╠╠╠╠╠╠╠▒        Ü▒╠╠╠╠╠╠╠╠╠▒
                 ╔╠╠╠╠╠╠╠╠╠╠╠R       [╠╠╠╠╠╠╠╠╠╠╠R       j╠╠╠╠╠╠╠╠╠╠╠R
                j╠╠╠╠╠╠╠╠╠╠╠╜       ╔╠╠╠╠╠╠╠╠╠╠╠R       ╔╠╠╠╠╠╠╠╠╠╠╠╩
               é╠╠╠╠╠╠╠╠╠╠╠Γ       j╠╠╠╠╠╠╠╠╠╠╠╜       ╔╠╠╠╠╠╠╠╠╠╠╠R
              ╬╠╠╠╠╠╠╠╠╠╠╠^       é╠╠╠╠╠╠╠╠╠╠╠Γ       j╠╠╠╠╠╠╠╠╠╠╠╜
             ╬╠╠╠╠╠╠╠╠╠╠╠`       ╬╠╠╠╠╠╠╠╠╠╠╠^       ê╠╠╠╠╠╠╠╠╠╠╠Γ
EOF
)

# Description of Devner Tools
devner_description=$(cat <<'EOF'

Devner
────────────
Devner is a comprehensive set of tools designed to streamline and optimize your development workflow. Whether you are managing databases, setting up project environments, or deploying developpement applications, Devner has you covered with its easy-to-use commands and efficient automation features.

Key Features:
- Rapid project setup for Laravel and WordPress.
- Seamless database management with support for MySQL and PostgreSQL.
- Integrated Docker environment management.
- Automated host configuration with Caddyfile management.
EOF
)

# Social Media Information
social_media=$(cat <<'EOF'

Connect with me
───────────────
GitHub: https://github.com/MarJC5
LinkedIn: https://linkedin.com/in/jean-christio-martin-385574111
Twitter: https://twitter.com/jeanchristio

EOF
)

# Display ASCII Art
echo "$ascii_art"

# Display Description of Devner Tools
echo "$devner_description"

# Display Social Media Information
echo "$social_media"
}