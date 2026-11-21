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
    private UUID id;

    @Column(nullable = false, unique = true, length = 32)
    private String code;

    @Column(nullable = false)
    private String name;

    @Column(nullable = false, length = 32)
    private String semester;

    @Column(name = "lead_instructor_id", nullable = false)
    private UUID leadInstructorId;

    @Column(name = "quota_account_ref", nullable = false, unique = true, length = 80)
    private String quotaAccountRef;

    @Column(name = "weekly_minutes", nullable = false)
    private int weeklyMinutes = 1200;

    @Column(name = "max_concurrent_sessions", nullable = false)
    private int maxConcurrentSessions = 30;

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

    public ClassEntity(String code, String name, String semester, UUID leadInstructorId, int weeklyMinutes, int maxConcurrentSessions) {
        this.id = UUID.randomUUID();
        this.code = code;
        this.name = name;
        this.semester = semester;
        this.leadInstructorId = leadInstructorId;
        this.quotaAccountRef = "class:" + this.id;
        this.weeklyMinutes = weeklyMinutes;
        this.maxConcurrentSessions = maxConcurrentSessions;
        this.status = ClassStatus.ACTIVE;
        this.createdAt = Instant.now();
        this.updatedAt = Instant.now();
    }

    public UUID getId() {
        return id;
    }

    public void setId(UUID id) {
        this.id = id;
        if (id != null) {
            this.quotaAccountRef = "class:" + id;
        }
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

    public UUID getLeadInstructorId() {
        return leadInstructorId;
    }

    public void setLeadInstructorId(UUID leadInstructorId) {
        this.leadInstructorId = leadInstructorId;
    }

    public String getQuotaAccountRef() {
        return quotaAccountRef;
    }

    public void setQuotaAccountRef(String quotaAccountRef) {
        this.quotaAccountRef = quotaAccountRef;
    }

    public int getWeeklyMinutes() {
        return weeklyMinutes;
    }

    public void setWeeklyMinutes(int weeklyMinutes) {
        this.weeklyMinutes = weeklyMinutes;
    }

    public int getMaxConcurrentSessions() {
        return maxConcurrentSessions;
    }

    public void setMaxConcurrentSessions(int maxConcurrentSessions) {
        this.maxConcurrentSessions = maxConcurrentSessions;
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
