cd ~/src/cart-ui-next
ISSUE_KEY=$(git branch | grep \* | cut -d ' ' -f2 |  awk -F'/' '{print $2}')
echo $ISSUE_KEY