-- SQLite disables foreign key enforcement by default
-- This setting applies per database connection rather than per database file
-- Must configure it right when the connection is initialized.
PRAGMA foreign_keys = ON;

INSERT INTO teachers (first_name, last_name, email, date_of_birth) VALUES
('John', 'Teacher', 'johnteacherz@gmail.com', '1995-10-10'),
('Jane', 'Teacher', 'janeteacherz@gmail.com', '1997-09-11'),
('Alex', 'Smith', 'alexsmith1092@gmail.com', '1990-08-22'),
('Antony', 'Herbert', 'tonyherbert10@gmail.com', '2000-01-23');

INSERT INTO students (first_name, last_name, email, date_of_birth) VALUES
('James', 'Millers', 'millerjs@gmail.com', '2008-03-14'),
('Christian', 'DeLuna', 'delunaaachris@gmail.com', '2008-07-22'),
('Justin', 'Milen', 'justinmilen88@gmail.com', '2007-11-05'),
('Cooper', 'Alexis', 'cooperlexisz@gmail.com', '2008-01-30'),
('Alexandra', 'McConnel', 'sandramcconnel09@gmail.com', '2008-03-21'),
('Jane', 'Grey', 'jngreyy777@gmail.com', '2007-03-12'),
('Casandra', 'Lopez', 'casandralopez223@gmail.com', '2006-12-28');

INSERT INTO subjects (name) VALUES
('Math'),
('English'),
('Science'),
('History'),
('Art');

INSERT INTO classes (name, homeroom_teacher_id) VALUES
('8A', 1),
('8B', 2),
('9A', NULL);

INSERT INTO teaching_assignments (class_id, subject_id, teacher_id) VALUES
(1, 1, 3), -- 8A - Math - Alex
(1, 2, 4), -- 8A - English - Antony
(1, 3, 2), -- 8A - Science - Jane
(2, 1, 3), -- 8B - Math - Alex
(2, 4, 1), -- 8B - History - John
(3, 3, 2), -- 9A - Science - Jane
(3, 5, 4); -- 9A - Art - Antony

INSERT INTO class_student_mapping (student_id, class_id) VALUES
(1, 1),
(2, 1),
(3, 2),
(4, 3);
