package dev.crucible.learning.domain;

import jakarta.persistence.*;
import java.time.Instant;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

@Entity
@Table(name = "classes")
public class ClassEntity {

    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    @Column(nullable = false, unique = true, length = 32)
    private String code;

    @Column(nullable = false)
    private String name;

    @Column(nullable = false, length = 32)
    private String semester;

    @Column(name = "instructor_quota_account_ref", length = 64)
    private String instructorQuotaAccountRef;

    @Column(name = "session_limit_per_student", nullable = false)
    private int sessionLimitPerStudent = 1;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false, length = 32)
    private ClassStatus status = ClassStatus.ACTIVE;

    @Column(name = "created_at", nullable = false, updatable = false)
    private Instant createdAt = Instant.now();

    @Column(name = "updated_at", nullable = false)
    private Instant updatedAt = Instant.now();

    @OneToMany(mappedBy = "clazz", cascade = CascadeType.ALL, orphanRemoval = true)
    private List<ClassMemberEntity> members = new ArrayList<>();

    public ClassEntity() {
    }

    public ClassEntity(String code, String name, String semester, String instructorQuotaAccountRef, int sessionLimitPerStudent) {
        this.code = code;
        this.name = name;
        this.semester = semester;
        this.instructorQuotaAccountRef = instructorQuotaAccountRef;
        this.sessionLimitPerStudent = sessionLimitPerStudent;
        this.status = ClassStatus.ACTIVE;
        this.createdAt = Instant.now();
        this.updatedAt = Instant.now();
    }

    public UUID getId() {
        return id;
    }

    public void setId(UUID id) {
        this.id = id;
    }

    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getSemester() {
        return semester;
    }

    public void setSemester(String semester) {
        this.semester = semester;
    }

    public String getInstructorQuotaAccountRef() {
        return instructorQuotaAccountRef;
    }

    public void setInstructorQuotaAccountRef(String instructorQuotaAccountRef) {
        this.instructorQuotaAccountRef = instructorQuotaAccountRef;
    }

    public int getSessionLimitPerStudent() {
        return sessionLimitPerStudent;
    }

    public void setSessionLimitPerStudent(int sessionLimitPerStudent) {
        this.sessionLimitPerStudent = sessionLimitPerStudent;
    }

    public ClassStatus getStatus() {
        return status;
    }

    public void setStatus(ClassStatus status) {
        this.status = status;
    }

    public Instant getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(Instant createdAt) {
        this.createdAt = createdAt;
    }

    public Instant getUpdatedAt() {
        return updatedAt;
    }

    public void setUpdatedAt(Instant updatedAt) {
        this.updatedAt = updatedAt;
    }

    public List<ClassMemberEntity> getMembers() {
        return members;
    }

    public void setMembers(List<ClassMemberEntity> members) {
        this.members = members;
    }

    public void addMember(ClassMemberEntity member) {
        members.add(member);
        member.setClazz(this);
    }
}
