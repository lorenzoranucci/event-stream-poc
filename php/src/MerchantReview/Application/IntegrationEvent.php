<?php

declare(strict_types=1);

namespace App\MerchantReview\Application;

abstract class IntegrationEvent
{
    private function __construct(
        private string $type,
        private string $message,
    ) {}

    abstract protected static function create(
        string $type,
        string $message,
    ): self;

    public function getType(): string
    {
        return $this->type;
    }

    public function getMessage(): string
    {
        return $this->message;
    }
}
