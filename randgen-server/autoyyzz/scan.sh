while read line
do
python template.py $line --tmpl tmpl_with_time
done < operator_list
