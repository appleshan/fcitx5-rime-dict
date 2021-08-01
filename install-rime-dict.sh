#!/usr/bin/env bash

INPUT_METHOD=fcitx5
RIME_CONFIG_PATH=~/.local/share/fcitx5/rime
DICT_DIR=$PWD

if grep -Eq "Arch|Manjaro" '/etc/os-release' || grep -Eq "Arch|Manjaro" '/etc/issue'; then
    LINUX_DISTRO='arch'
fi

if [ "${LINUX_DISTRO}" != "arch" ]; then
    echo '截至20210801，本功能暂只适配Arch系发行版'
    exit
fi

if ! command -v ${INPUT_METHOD} &>/dev/null; then
    echo "${INPUT_METHOD} could not be found"
    echo "Please install ${INPUT_METHOD} in your distro and try again!"
    exit
fi

# backup rime config
# if [ -d ${RIME_CONFIG_PATH} ]; then
#     BACKUP_PATH=$(mktemp -p $(dirname ${RIME_CONFIG_PATH}) -d \
#         -t ${RIMERC}-backup-$(date '+%Y%m%d')-XXXX)
#     mv ${RIME_CONFIG_PATH}/* ${BACKUP_PATH}
#     echo -e "\nBackup your rime config to ${BACKUP_PATH}"
# fi

mkdir -p ${RIME_CONFIG_PATH}

# copy config to config path
# cp -ra ./* ${RIME_CONFIG_PATH}
ln -sf "${DICT_DIR}"/default.custom.yaml            "${RIME_CONFIG_PATH}"/default.custom.yaml
ln -sf "${DICT_DIR}"/luna_pinyin.custom.yaml        "${RIME_CONFIG_PATH}"/luna_pinyin.custom.yaml
ln -sf "${DICT_DIR}"/luna_pinyin.extended.dict.yaml "${RIME_CONFIG_PATH}"/luna_pinyin.extended.dict.yaml
ln -sf "${DICT_DIR}"/luna_pinyin_simp.custom.yaml   "${RIME_CONFIG_PATH}"/luna_pinyin_simp.custom.yaml
ln -sf "${DICT_DIR}"/easy_en.dict.yaml              "${RIME_CONFIG_PATH}"/easy_en.dict.yaml
ln -sf "${DICT_DIR}"/easy_en.schema.yaml            "${RIME_CONFIG_PATH}"/easy_en.schema.yaml
ln -sf "${DICT_DIR}"/easy_en.yaml                   "${RIME_CONFIG_PATH}"/easy_en.yaml
ln -sf "${DICT_DIR}"/numbers.schema.yaml            "${RIME_CONFIG_PATH}"/numbers.schema.yaml
ln -sf "${DICT_DIR}"/symbols.dict.yaml              "${RIME_CONFIG_PATH}"/symbols.dict.yaml
ln -sf "${DICT_DIR}"/symbols.schema.yaml            "${RIME_CONFIG_PATH}"/symbols.schema.yaml
ln -sf "${DICT_DIR}"/custom_phrase.txt              "${RIME_CONFIG_PATH}"/custom_phrase.txt

ln -sf "${DICT_DIR}"/lua/easy_en.lua "${RIME_CONFIG_PATH}"/lua/easy_en.lua

ln -sf "${DICT_DIR}"/luna_pinyin.extended/corpuscharacterlist.dict.yaml "${RIME_CONFIG_PATH}"/luna_pinyin.extended/corpuscharacterlist.dict.yaml
ln -sf "${DICT_DIR}"/luna_pinyin.extended/corpuswordlist.dict.yaml      "${RIME_CONFIG_PATH}"/luna_pinyin.extended/corpuswordlist.dict.yaml
ln -sf "${DICT_DIR}"/luna_pinyin.extended/sogouw.dict.yaml              "${RIME_CONFIG_PATH}"/luna_pinyin.extended/sogouw.dict.yaml
ln -sf "${DICT_DIR}"/luna_pinyin.extended/polyphones.dict.yaml          "${RIME_CONFIG_PATH}"/luna_pinyin.extended/polyphones.dict.yaml
ln -sf "${DICT_DIR}"/luna_pinyin.extended/symbols.dict.yaml             "${RIME_CONFIG_PATH}"/luna_pinyin.extended/symbols.dict.yaml
ln -sf "${DICT_DIR}"/luna_pinyin.extended/zhwiki-20210722.dict.yaml     "${RIME_CONFIG_PATH}"/luna_pinyin.extended/zhwiki-20210722.dict.yaml
ln -sf "${DICT_DIR}"/luna_pinyin.extended/zhwiki-20210722.dict          "${RIME_CONFIG_PATH}"/luna_pinyin.extended/zhwiki-20210722.dict
ln -sf "${DICT_DIR}"/luna_pinyin.extended/essay.dict.yaml               "${RIME_CONFIG_PATH}"/luna_pinyin.extended/essay.dict.yaml
ln -sf "${DICT_DIR}"/luna_pinyin.extended/sogou_chengyusuyu.dict.yaml   "${RIME_CONFIG_PATH}"/luna_pinyin.extended/sogou_chengyusuyu.dict.yaml
ln -sf "${DICT_DIR}"/luna_pinyin.extended/sogou_shici.dict.yaml         "${RIME_CONFIG_PATH}"/luna_pinyin.extended/sogou_shici.dict.yaml
ln -sf "${DICT_DIR}"/luna_pinyin.extended/sogou_yaoming.dict.yaml       "${RIME_CONFIG_PATH}"/luna_pinyin.extended/sogou_yaoming.dict.yaml
ln -sf "${DICT_DIR}"/luna_pinyin.extended/cn_en.dict.yaml               "${RIME_CONFIG_PATH}"/luna_pinyin.extended/cn_en.dict.yaml
ln -sf "${DICT_DIR}"/luna_pinyin.extended/kaifa.dict.yaml               "${RIME_CONFIG_PATH}"/luna_pinyin.extended/kaifa.dict.yaml
ln -sf "${DICT_DIR}"/luna_pinyin.extended/food.dict.yaml                "${RIME_CONFIG_PATH}"/luna_pinyin.extended/food.dict.yaml

# copy userdb to config path
# cp -ra ${BACKUP_PATH}/*.userdb ${RIME_CONFIG_PATH}
# cp -ra ${BACKUP_PATH}/user.yaml ${RIME_CONFIG_PATH}
# cp -ra ${BACKUP_PATH}/installation.yaml ${RIME_CONFIG_PATH}

# restart input method
echo -e "\nWaiting for ${INPUT_METHOD} restart and deploying..."

nohup ${INPUT_METHOD} -r >/dev/null 2>&1 &

echo -e "\nSetup fcitx5-rime-dict successfully!"
