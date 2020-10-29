# Usage: bin/create_service.sh foo-service

set -e

mkdir "./services/${1}"

# -n Do not overwrite an existing file
# -R Recurse into directory. Source must end in slash
#    to copy the contents of the directory
cp -nR ./tools/templates/service/ "./services/${1}/"

# https://stackoverflow.com/a/50224830
service_name_pascal=$(echo "${1}" | perl -nE 'say join "", map {ucfirst lc} split /[^[:alnum:]]+/')
service_name_snake=$(echo "${1}" | tr '-' '_')

# Loop through all the newly copied files & directories
files=$(find "./services/${1}")
for f in $files; do
    # Replace any instances of {{service_name_snake}} in the file (or directory) names
    new_name=$(echo $f | sed "s/{{service_name_snake}}/${service_name_snake}/")

    # Remove the .tmpl extensions from everything
    # https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html#Shell-Parameter-Expansion
    new_name=${new_name%.tmpl}

    if [ "$f" != "$new_name" ]; then
        mv -n "$f" "$new_name"
    fi
done

# Loop through all of the files only and replace strings inside
files=$(find "./services/${1}" -type f)
for f in $files; do
    sed -i '' "s/{{service_name_kebab}}/${1}/" "$f"
    sed -i '' "s/{{service_name_pascal}}/${service_name_pascal}/" "$f"
    sed -i '' "s/{{service_name_snake}}/${service_name_snake}/" "$f"
done

go generate "./services/${1}/..."
