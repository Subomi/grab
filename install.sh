#!/bin/sh
#TODO: Make this script as cross compatible as possible.
#TODO: Read this through and parse with grep or sed or awk. 
#TODO: Purge previous installation and create new one.
USERHOME=`printenv HOME`

DEFAULT_FOLDER_NAME="/.grabber"
DEFAULT_SYSTEM_FOLDER="/etc/grabber"
DEFAULT_SYSTEM_CONF="./conf/system.yml"
DEFAULT_USER_CONF="./conf/user.yml"

echo "Installing grab .."

DEFAULT_CACHE_DIR="$USERHOME/cache"
echo $DEFAULT_CACHE_DIR
GRABBERHOME="$USERHOME$DEFAULT_FOLDER_FOLDER"

check_and_create_folder() 
{
    if [ -d $1 ]
    then   
        echo "$1 exists already."
    else
        echo "Creating $1 .."
        mkdir $1
        if [ "$?" -ne "0" ]
        then 
            echo "An error occurred creating $1 folder"
        fi
    fi
}

copy_default_files() 
{
    [ -f $1 ] && echo "" || echo "$1 doesn't exist, your download doesn't seem complete."
    cp -v -R $1 $2
    if [ "$?" -ne "0" ]
    then
        echo "An error occurred copying default conf files."
    fi 
}

# Make Folders
check_and_create_folder $GRABBERHOME
check_and_create_folder $DEFAULT_SYSTEM_FOLDER
check_and_create_folder $DEFAULT_CACHE_DIR

# Copy files over
copy_default_files $DEFAULT_SYSTEM_CONF $DEFAULT_SYSTEM_FOLDER
copy_default_files $DEFAULT_USER_CONF $GRABBERHOME
