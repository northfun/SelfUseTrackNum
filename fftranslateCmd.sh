#! /bin/sh

checkInt1(){ 
    tmp=`echo $1 |sed 's/[0-9]//g'` 
    # [ -n "${tmp}" ] && { echo $1"Args must be integer!";echo 1; } 
    [ -n "${tmp}" ] && echo 1
}

#getparams(){
#    echo $1
#    list=`$1`
#    for var in $list ; do
#        res=`checkInt1 $var`
#        if [ "$res" == "1" ];then
#            continue
#        fi
#        if [ "$params" == "" ]; then
#            params=$var
#        else
#            params=$params","$var
#        fi
#    done
#    echo $params
#}

diffcmd(){
    params=""
    list=`git diff | grep USER.*$1.*PARAM.*= | cut -d'=' -f 2 `
    for var in ${list} ; do
        res=`checkInt1 $var`
        if [ "$res" == "1" ];then
            continue
        fi
        if [ "$params" == "" ]; then
            params=$var
        else
            params=$params","$var
        fi
    done
    echo $params
}

allcmd(){
    params=""
    list=` grep USER.*$1.*PARAM.*= $2 | cut -d'=' -f 2 `
    for var in ${list} ; do
        res=`checkInt1 $var`
        if [ "$res" == "1" ];then
            continue
        fi
        if [ "$params" == "" ]; then
            params=$var
        else
            params=$params","$var
        fi
    done
    echo $params
}

case $1 in
	allcmd)
		allcmd $2 $3
	;;
    diffcmd)
        diffcmd $2
    ;;
	*)
	    echo default	
	;;
esac 
