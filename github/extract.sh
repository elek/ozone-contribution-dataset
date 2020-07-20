echo "number,title,login" > participant.csv
jq -r '. as $pr | .participants.edges[].node | [$pr.number, $pr.title, .login] | @csv' prs/*.json > participant.csv

echo "number,title,login,state,updated" > reviews.csv
jq -r '. as $pr | .reviews.nodes[] | [$pr.number, $pr.title, .author.login, .state, .updatedAt] | @csv' prs/*.json >> reviews.csv

echo "number,title, state,merged,createdAt,updatedAt,mergedAt,closedAt,closed,baseRefName,author" > prs.csv
jq -r '[.number, .title, .state, .merged, .createdAt, .updatedAt, .mergedAt, .closedAt, .closed, .baseRefName,.author.login] | @csv' prs/*.json >> prs.csv

echo "number,title,author,role,createdAt" > comments.csv
jq -r '. as $pr | .comments.nodes[] | [$pr.number, $pr.title, .author.login, .authorAssociation, .createdAt] | @csv' prs/*.json >>comments.csv
