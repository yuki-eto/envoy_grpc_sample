#!/usr/bin/env bash

set -e

OLD_443_ACTIONS=$(aws elbv2 describe-listeners --listener-arns ${LISTENER_443_ARN} | jq -c ".Listeners[0].DefaultActions")
OLD_10443_ACTIONS=$(aws elbv2 describe-listeners --listener-arns ${LISTENER_10443_ARN} | jq -c ".Listeners[0].DefaultActions")

echo "- change 443 listener action"
aws elbv2 modify-listener --listener-arn ${LISTENER_443_ARN} --default-actions ${OLD_10443_ACTIONS}
echo "- change 10443 listener action"
aws elbv2 modify-listener --listener-arn ${LISTENER_10443_ARN} --default-actions ${OLD_443_ACTIONS}
