Implementation of the very opinionated RESTful specification used internally at Khaos Systems, and a bunch of ther Khaos Group companies.

Repository Pattern: https://www.youtube.com/watch?v=ivJ2s0e7vi0

Guidelines:
 - Pointer fields are for relations. Reason: If a field is stored by value, it's implicitly a part of the entity, and thus can not be reference by another entity. It builds an implicit ownership relationship.
 - Value fields are embedded.

TODO:
 - Make create use schema to fetch fields.
 - Support omitempty for fields, and if fields are expanded, include them in the response even if they are empty.
 - Recursive expand.
 - DeletedAt field for soft deletes. (could be tags on fields, to make it optional)
 - CreatedAt and UpdatedAt fields. (could be tags on fields, to make it optional)