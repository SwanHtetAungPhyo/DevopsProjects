find . -name "*.go" -type f | while read file; do
    echo "Updating $file"
    sed -i '' 's|"awesomeProject|"github.com/SwanHtetAungPhyo/docker-log-agent|g' "$file"
done