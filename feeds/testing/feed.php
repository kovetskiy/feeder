<?php

# only for development.
# will be removed when branch will be merged with master

$entries = [
    [
        'uid' => md5('some-uid' . date('d.m.Y H')),
        'title'=> 'Valid entry',
        'url'=> 'http=>//blahblah111',
        'preview' => 'Ohhh that is preview',
        'create_date' => time(),
    ],
    [
        'uid' => 'Invalid entry',
        'title'=> 'Invalid entry',
    ],
];

print json_encode($entries);
