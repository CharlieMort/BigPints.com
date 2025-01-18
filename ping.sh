for ip in $(seq 0 255); do

ping -c 1 192.168.0.$ip | grep "Reply" &

done